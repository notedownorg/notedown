// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package extensions

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestWikilinkRegex(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantMatch   bool
		wantTarget  string
		wantDisplay string
	}{
		{
			name:        "simple wikilink",
			input:       "[[page]]",
			wantMatch:   true,
			wantTarget:  "page",
			wantDisplay: "",
		},
		{
			name:        "wikilink with display text",
			input:       "[[target|display]]",
			wantMatch:   true,
			wantTarget:  "target",
			wantDisplay: "display",
		},
		{
			name:        "wikilink with spaces",
			input:       "[[my page|My Page]]",
			wantMatch:   true,
			wantTarget:  "my page",
			wantDisplay: "My Page",
		},
		{
			name:        "wikilink at start of line",
			input:       "[[start]] of line",
			wantMatch:   true,
			wantTarget:  "start",
			wantDisplay: "",
		},
		{
			name:      "incomplete wikilink - single bracket",
			input:     "[page]",
			wantMatch: false,
		},
		{
			name:      "incomplete wikilink - missing closing",
			input:     "[[page",
			wantMatch: false,
		},
		{
			name:      "empty wikilink",
			input:     "[[]]",
			wantMatch: false,
		},
		{
			name:      "wikilink with empty display",
			input:     "[[page|]]",
			wantMatch: false, // Empty display is not valid
		},
		{
			name:      "wikilink with only pipe",
			input:     "[[|]]",
			wantMatch: false,
		},
		{
			name:      "nested brackets in target",
			input:     "[[page]with]brackets]]",
			wantMatch: false,
		},
		{
			name:        "relative path with extension",
			input:       "[[docs/api.md]]",
			wantMatch:   true,
			wantTarget:  "docs/api.md",
			wantDisplay: "",
		},
		{
			name:        "relative path without extension",
			input:       "[[foo/bar]]",
			wantMatch:   true,
			wantTarget:  "foo/bar",
			wantDisplay: "",
		},
		{
			name:        "relative path with display text",
			input:       "[[foo/bar.md|Foo Bar]]",
			wantMatch:   true,
			wantTarget:  "foo/bar.md",
			wantDisplay: "Foo Bar",
		},
		{
			name:        "deep nested path",
			input:       "[[projects/alpha/docs/readme.md]]",
			wantMatch:   true,
			wantTarget:  "projects/alpha/docs/readme.md",
			wantDisplay: "",
		},
		{
			name:        "path with spaces and display",
			input:       "[[docs/user guide.md|User Guide]]",
			wantMatch:   true,
			wantTarget:  "docs/user guide.md",
			wantDisplay: "User Guide",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := wikilinkRegex.FindSubmatch([]byte(tt.input))

			if tt.wantMatch {
				if match == nil {
					t.Errorf("Expected match for input %q, got nil", tt.input)
					return
				}

				target := string(match[1])
				if target != tt.wantTarget {
					t.Errorf("Expected target %q, got %q", tt.wantTarget, target)
				}

				var display string
				if len(match) > 2 && match[2] != nil {
					display = string(match[2])
				}
				if display != tt.wantDisplay {
					t.Errorf("Expected display %q, got %q", tt.wantDisplay, display)
				}
			} else {
				if match != nil {
					t.Errorf("Expected no match for input %q, got %v", tt.input, match)
				}
			}
		})
	}
}

func TestWikilinkParser_Trigger(t *testing.T) {
	parser := &wikilinkParser{}
	triggers := parser.Trigger()

	if len(triggers) != 1 || triggers[0] != '[' {
		t.Errorf("Expected trigger '[', got %v", triggers)
	}
}

func TestWikilinkParser_Parse(t *testing.T) {
	// Test basic parsing functionality with integration tests
	// This avoids the complex goldmark parser context API issues
	tests := []struct {
		name        string
		markdown    string
		wantTargets []string
		wantDisplay []string
	}{
		{
			name:        "simple wikilink",
			markdown:    "Text with [[page]] link",
			wantTargets: []string{"page"},
			wantDisplay: []string{"page"},
		},
		{
			name:        "wikilink with display text",
			markdown:    "Text with [[target|display]] link",
			wantTargets: []string{"target"},
			wantDisplay: []string{"display"},
		},
		{
			name:        "multiple wikilinks",
			markdown:    "[[link1]] and [[link2|Link 2]]",
			wantTargets: []string{"link1", "link2"},
			wantDisplay: []string{"link1", "Link 2"},
		},
		{
			name:        "relative paths",
			markdown:    "See [[docs/api.md]] and [[foo/bar]]",
			wantTargets: []string{"docs/api.md", "foo/bar"},
			wantDisplay: []string{"docs/api.md", "foo/bar"},
		},
		{
			name:        "relative paths with display text",
			markdown:    "Check [[docs/api.md|API Docs]] and [[projects/alpha|Alpha Project]]",
			wantTargets: []string{"docs/api.md", "projects/alpha"},
			wantDisplay: []string{"API Docs", "Alpha Project"},
		},
		{
			name:        "no wikilinks",
			markdown:    "Just regular text [not a link]",
			wantTargets: []string{},
			wantDisplay: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(NewWikilinkExtension()))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var wikilinks []*WikilinkAST
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if wl, ok := node.(*WikilinkAST); ok {
						wikilinks = append(wikilinks, wl)
					}
				}
				return ast.WalkContinue, nil
			})

			if len(wikilinks) != len(tt.wantTargets) {
				t.Errorf("Expected %d wikilinks, got %d", len(tt.wantTargets), len(wikilinks))
				return
			}

			for i, wl := range wikilinks {
				if wl.Target != tt.wantTargets[i] {
					t.Errorf("Wikilink %d: expected target %q, got %q", i, tt.wantTargets[i], wl.Target)
				}
				if wl.DisplayText != tt.wantDisplay[i] {
					t.Errorf("Wikilink %d: expected display %q, got %q", i, tt.wantDisplay[i], wl.DisplayText)
				}
			}
		})
	}
}

func TestWikilinkAST_Kind(t *testing.T) {
	node := &WikilinkAST{
		Target:      "test",
		DisplayText: "Test",
	}

	if node.Kind() != WikilinkKind {
		t.Errorf("Expected kind %v, got %v", WikilinkKind, node.Kind())
	}
}

func TestWikilinkAST_Dump(t *testing.T) {
	node := &WikilinkAST{
		Target:      "test-target",
		DisplayText: "Test Display",
	}

	// Test that Dump doesn't panic - we can't easily test the output
	// since it writes to internal goldmark structures
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump panicked: %v", r)
		}
	}()

	node.Dump([]byte("test source"), 0)
}

func TestWikilinkKind(t *testing.T) {
	// Test that WikilinkKind is properly initialized
	if WikilinkKind.String() != "Wikilink" {
		t.Errorf("Expected WikilinkKind to be 'Wikilink', got %q", WikilinkKind.String())
	}
}

func TestWikilinkIntegration(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantHTML string
	}{
		{
			name:     "simple wikilink",
			markdown: "Text with [[wikilink]] in it.",
			// Note: Without a renderer, wikilinks won't render as HTML
			// This test verifies the AST structure is correct
		},
		{
			name:     "wikilink with display text",
			markdown: "Text with [[target|display text]] in it.",
		},
		{
			name:     "multiple wikilinks",
			markdown: "First [[link1]] and second [[link2|Link 2]].",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(NewWikilinkExtension()))

			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			// Walk the AST to find wikilink nodes
			var wikilinks []*WikilinkAST
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if wl, ok := node.(*WikilinkAST); ok {
						wikilinks = append(wikilinks, wl)
					}
				}
				return ast.WalkContinue, nil
			})

			// For simple wikilink test, expect 1 wikilink
			if tt.name == "simple wikilink" {
				if len(wikilinks) != 1 {
					t.Errorf("Expected 1 wikilink, got %d", len(wikilinks))
					return
				}
				if wikilinks[0].Target != "wikilink" {
					t.Errorf("Expected target 'wikilink', got %q", wikilinks[0].Target)
				}
				if wikilinks[0].DisplayText != "wikilink" {
					t.Errorf("Expected display 'wikilink', got %q", wikilinks[0].DisplayText)
				}
			}

			// For wikilink with display text test
			if tt.name == "wikilink with display text" {
				if len(wikilinks) != 1 {
					t.Errorf("Expected 1 wikilink, got %d", len(wikilinks))
					return
				}
				if wikilinks[0].Target != "target" {
					t.Errorf("Expected target 'target', got %q", wikilinks[0].Target)
				}
				if wikilinks[0].DisplayText != "display text" {
					t.Errorf("Expected display 'display text', got %q", wikilinks[0].DisplayText)
				}
			}

			// For multiple wikilinks test
			if tt.name == "multiple wikilinks" {
				if len(wikilinks) != 2 {
					t.Errorf("Expected 2 wikilinks, got %d", len(wikilinks))
					return
				}

				// Check first wikilink
				if wikilinks[0].Target != "link1" {
					t.Errorf("Expected first target 'link1', got %q", wikilinks[0].Target)
				}
				if wikilinks[0].DisplayText != "link1" {
					t.Errorf("Expected first display 'link1', got %q", wikilinks[0].DisplayText)
				}

				// Check second wikilink
				if wikilinks[1].Target != "link2" {
					t.Errorf("Expected second target 'link2', got %q", wikilinks[1].Target)
				}
				if wikilinks[1].DisplayText != "Link 2" {
					t.Errorf("Expected second display 'Link 2', got %q", wikilinks[1].DisplayText)
				}
			}
		})
	}
}

func TestWikilinkContextualParsing(t *testing.T) {
	// Test that wikilinks are parsed correctly in different markdown contexts
	tests := []struct {
		name      string
		markdown  string
		wantCount int
	}{
		{
			name:      "wikilink in paragraph",
			markdown:  "This is a paragraph with [[wikilink]].",
			wantCount: 1,
		},
		{
			name:      "wikilink in heading",
			markdown:  "# Heading with [[wikilink]]",
			wantCount: 1,
		},
		{
			name:      "wikilink in list item",
			markdown:  "- List item with [[wikilink]]",
			wantCount: 1,
		},
		{
			name:      "wikilink in emphasis",
			markdown:  "*Emphasized [[wikilink]] text*",
			wantCount: 1,
		},
		{
			name:      "wikilink in blockquote",
			markdown:  "> Quote with [[wikilink]]",
			wantCount: 1,
		},
		{
			name:      "no wikilink in code block",
			markdown:  "```\n[[not-a-wikilink]]\n```",
			wantCount: 0,
		},
		{
			name:      "no wikilink in inline code",
			markdown:  "Code: `[[not-a-wikilink]]`",
			wantCount: 0,
		},
		{
			name:      "relative path in paragraph",
			markdown:  "This paragraph links to [[docs/api.md]].",
			wantCount: 1,
		},
		{
			name:      "relative path in heading",
			markdown:  "# API Reference [[docs/api.md]]",
			wantCount: 1,
		},
		{
			name:      "relative path in list item",
			markdown:  "- Check [[projects/alpha/readme.md]] for details",
			wantCount: 1,
		},
		{
			name:      "multiple relative paths",
			markdown:  "See [[docs/api.md]] and [[foo/bar]] and [[projects/alpha.md|Alpha]]",
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(NewWikilinkExtension()))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var count int
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if _, ok := node.(*WikilinkAST); ok {
						count++
					}
				}
				return ast.WalkContinue, nil
			})

			if count != tt.wantCount {
				t.Errorf("Expected %d wikilinks, got %d", tt.wantCount, count)
			}
		})
	}
}

func TestWikilinkEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMatch bool
	}{
		{
			name:      "unicode in target",
			input:     "[[café]]",
			wantMatch: true,
		},
		{
			name:      "numbers in target",
			input:     "[[page123]]",
			wantMatch: true,
		},
		{
			name:      "special chars in target",
			input:     "[[page-with_special.chars]]",
			wantMatch: true,
		},
		{
			name:      "very long target",
			input:     "[[" + string(make([]byte, 1000)) + "]]",
			wantMatch: true,
		},
		{
			name:      "multiple pipes",
			input:     "[[target|display|extra]]",
			wantMatch: true, // regex will match but only use first pipe
		},
		{
			name:      "relative path",
			input:     "[[docs/api.md]]",
			wantMatch: true,
		},
		{
			name:      "nested path with extension",
			input:     "[[projects/alpha/readme.md]]",
			wantMatch: true,
		},
		{
			name:      "relative path with display text",
			input:     "[[foo/bar.md|Foo Bar]]",
			wantMatch: true,
		},
		{
			name:      "deep nested path",
			input:     "[[docs/technical/architecture/database.md]]",
			wantMatch: true,
		},
		{
			name:      "relative path without extension",
			input:     "[[foo/bar]]",
			wantMatch: true,
		},
		{
			name:      "nested path without extension",
			input:     "[[docs/api/endpoints]]",
			wantMatch: true,
		},
		{
			name:      "path without extension with display",
			input:     "[[projects/alpha|Project Alpha]]",
			wantMatch: true,
		},
		{
			name:      "directory traversal attempt - should be rejected",
			input:     "[[../parent-doc]]",
			wantMatch: false,
		},
		{
			name:      "directory traversal in subdirectory - should be rejected",
			input:     "[[docs/../secret]]",
			wantMatch: false,
		},
		{
			name:      "multiple directory traversal - should be rejected",
			input:     "[[../../outside-workspace]]",
			wantMatch: false,
		},
		{
			name:      "directory traversal with display text - should be rejected",
			input:     "[[../parent-doc|Parent Document]]",
			wantMatch: false,
		},
	}

	md := goldmark.New(goldmark.WithExtensions(NewWikilinkExtension()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := md.Parser().Parse(text.NewReader([]byte(tt.input)))

			var found bool
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if _, ok := node.(*WikilinkAST); ok {
						found = true
					}
				}
				return ast.WalkContinue, nil
			})

			if found != tt.wantMatch {
				t.Errorf("Expected wikilink found = %v, got %v", tt.wantMatch, found)
			}
		})
	}
}
