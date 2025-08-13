package notedownls

import (
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

func TestGetWikilinkContext(t *testing.T) {
	server := &Server{
		logger: log.NewDefault(),
	}

	tests := []struct {
		name     string
		content  string
		position lsp.Position
		want     *WikilinkContext
	}{
		{
			name:     "no wikilink",
			content:  "Some regular text",
			position: lsp.Position{Line: 0, Character: 5},
			want:     nil,
		},
		{
			name:     "cursor before wikilink",
			content:  "Some [[link]] text",
			position: lsp.Position{Line: 0, Character: 3},
			want:     nil,
		},
		{
			name:     "cursor at start of wikilink",
			content:  "Some [[link]] text",
			position: lsp.Position{Line: 0, Character: 7},
			want: &WikilinkContext{
				Prefix:     "",
				IsComplete: true,
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 7},
					End:   lsp.Position{Line: 0, Character: 7},
				},
			},
		},
		{
			name:     "cursor in middle of wikilink",
			content:  "Some [[link]] text",
			position: lsp.Position{Line: 0, Character: 9},
			want: &WikilinkContext{
				Prefix:     "li",
				IsComplete: true,
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 7},
					End:   lsp.Position{Line: 0, Character: 9},
				},
			},
		},
		{
			name:     "cursor at end of incomplete wikilink",
			content:  "Some [[link",
			position: lsp.Position{Line: 0, Character: 11},
			want: &WikilinkContext{
				Prefix:     "link",
				IsComplete: false,
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 7},
					End:   lsp.Position{Line: 0, Character: 11},
				},
			},
		},
		{
			name:     "cursor in incomplete wikilink with partial text",
			content:  "Some [[proj",
			position: lsp.Position{Line: 0, Character: 11},
			want: &WikilinkContext{
				Prefix:     "proj",
				IsComplete: false,
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 7},
					End:   lsp.Position{Line: 0, Character: 11},
				},
			},
		},
		{
			name:     "multiple wikilinks, cursor in second",
			content:  "See [[first]] and [[second",
			position: lsp.Position{Line: 0, Character: 26},
			want: &WikilinkContext{
				Prefix:     "second",
				IsComplete: false,
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 20},
					End:   lsp.Position{Line: 0, Character: 26},
				},
			},
		},
		{
			name:     "wikilink with pipe separator",
			content:  "See [[target|display]] text",
			position: lsp.Position{Line: 0, Character: 12},
			want: &WikilinkContext{
				Prefix:     "target",
				IsComplete: true,
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 6},
					End:   lsp.Position{Line: 0, Character: 12},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{
				URI:      "file:///test.md",
				Basepath: "test",
				Content:  tt.content,
				Version:  1,
			}
			got := server.getWikilinkContext(doc, tt.position)

			if tt.want == nil && got != nil {
				t.Errorf("expected nil, got %+v", got)
				return
			}
			if tt.want != nil && got == nil {
				t.Errorf("expected %+v, got nil", tt.want)
				return
			}
			if tt.want == nil && got == nil {
				return
			}

			if got.Prefix != tt.want.Prefix {
				t.Errorf("prefix mismatch: got %q, want %q", got.Prefix, tt.want.Prefix)
			}
			if got.IsComplete != tt.want.IsComplete {
				t.Errorf("isComplete mismatch: got %v, want %v", got.IsComplete, tt.want.IsComplete)
			}
			if got.Range.Start.Line != tt.want.Range.Start.Line ||
				got.Range.Start.Character != tt.want.Range.Start.Character ||
				got.Range.End.Line != tt.want.Range.End.Line ||
				got.Range.End.Character != tt.want.Range.End.Character {
				t.Errorf("range mismatch: got %+v, want %+v", got.Range, tt.want.Range)
			}
		})
	}
}

func TestGenerateWikilinkTargets(t *testing.T) {
	server := &Server{
		logger: log.NewDefault(),
	}

	tests := []struct {
		name     string
		fileInfo *FileInfo
		want     []WikilinkTarget
	}{
		{
			name: "simple file in root",
			fileInfo: &FileInfo{
				URI:  "file:///workspace/readme.md",
				Path: "readme.md",
			},
			want: []WikilinkTarget{
				{
					Link:    "readme",
					Detail:  "Link to readme.md",
					SortKey: "0_readme",
				},
			},
		},
		{
			name: "file in subdirectory",
			fileInfo: &FileInfo{
				URI:  "file:///workspace/docs/api.md",
				Path: "docs/api.md",
			},
			want: []WikilinkTarget{
				{
					Link:    "api",
					Detail:  "Link to docs/api.md",
					SortKey: "0_api",
				},
				{
					Link:    "docs/api",
					Detail:  "Link to docs/api.md",
					SortKey: "1_docs/api",
				},
			},
		},
		{
			name: "file in nested subdirectory",
			fileInfo: &FileInfo{
				URI:  "file:///workspace/projects/alpha/readme.md",
				Path: "projects/alpha/readme.md",
			},
			want: []WikilinkTarget{
				{
					Link:    "readme",
					Detail:  "Link to projects/alpha/readme.md",
					SortKey: "0_readme",
				},
				{
					Link:    "projects/alpha/readme",
					Detail:  "Link to projects/alpha/readme.md",
					SortKey: "1_projects/alpha/readme",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server.generateWikilinkTargets(tt.fileInfo, "")

			if len(got) != len(tt.want) {
				t.Errorf("length mismatch: got %d targets, want %d", len(got), len(tt.want))
				return
			}

			for i, target := range got {
				want := tt.want[i]
				if target.Link != want.Link {
					t.Errorf("target[%d].Link = %q, want %q", i, target.Link, want.Link)
				}
				if target.Detail != want.Detail {
					t.Errorf("target[%d].Detail = %q, want %q", i, target.Detail, want.Detail)
				}
				if target.SortKey != want.SortKey {
					t.Errorf("target[%d].SortKey = %q, want %q", i, target.SortKey, want.SortKey)
				}
			}
		})
	}
}

func TestGetWikilinkCompletions(t *testing.T) {
	server := &Server{
		logger: log.NewDefault(),
		workspace: &WorkspaceManager{
			fileIndex: map[string]*FileInfo{
				"file:///workspace/readme.md": {
					URI:  "file:///workspace/readme.md",
					Path: "readme.md",
				},
				"file:///workspace/docs/api.md": {
					URI:  "file:///workspace/docs/api.md",
					Path: "docs/api.md",
				},
				"file:///workspace/projects/alpha.md": {
					URI:  "file:///workspace/projects/alpha.md",
					Path: "projects/alpha.md",
				},
				"file:///workspace/projects/beta.md": {
					URI:  "file:///workspace/projects/beta.md",
					Path: "projects/beta.md",
				},
			},
		},
	}

	tests := []struct {
		name         string
		prefix       string
		currentURI   string
		wantCount    int
		wantContains []string
	}{
		{
			name:         "empty prefix returns all files except current",
			prefix:       "",
			currentURI:   "file:///workspace/readme.md",
			wantCount:    6, // 3 files * 2 targets each (base name + path), minus current file
			wantContains: []string{"api", "docs/api", "alpha", "projects/alpha"},
		},
		{
			name:         "prefix 'a' filters correctly",
			prefix:       "a",
			currentURI:   "file:///workspace/readme.md",
			wantCount:    2, // "api" and "alpha"
			wantContains: []string{"api", "alpha"},
		},
		{
			name:         "prefix 'proj' filters to project files",
			prefix:       "proj",
			currentURI:   "file:///workspace/readme.md",
			wantCount:    2, // "projects/alpha" and "projects/beta"
			wantContains: []string{"projects/alpha", "projects/beta"},
		},
		{
			name:         "prefix matching nothing returns empty",
			prefix:       "xyz",
			currentURI:   "file:///workspace/readme.md",
			wantCount:    0,
			wantContains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server.getWikilinkCompletions(tt.prefix, tt.currentURI)

			if len(got) != tt.wantCount {
				t.Errorf("got %d completions, want %d", len(got), tt.wantCount)
				for _, item := range got {
					t.Logf("  - %s", item.Label)
				}
			}

			// Check that all expected items are present
			foundItems := make(map[string]bool)
			for _, item := range got {
				foundItems[item.Label] = true
			}

			for _, expected := range tt.wantContains {
				if !foundItems[expected] {
					t.Errorf("expected to find completion item %q but didn't", expected)
				}
			}
		})
	}
}
