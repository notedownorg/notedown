package notedownls

import (
	"testing"
)

func TestNewDocument(t *testing.T) {
	tests := []struct {
		name        string
		uri         string
		wantPath    string
		shouldError bool
	}{
		{
			name:     "file URI",
			uri:      "file:///path/to/document.md",
			wantPath: "document",
		},
		{
			name:     "file URI with relative path",
			uri:      "file://./README.md",
			wantPath: "README",
		},
		{
			name:     "non-file URI",
			uri:      "http://example.com/doc",
			wantPath: "",
		},
		{
			name:     "file with multiple extensions",
			uri:      "file:///path/to/archive.tar.gz",
			wantPath: "archive.tar",
		},
		{
			name:     "file without extension",
			uri:      "file:///path/to/LICENSE",
			wantPath: "LICENSE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewDocument(tt.uri)
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if doc.URI != tt.uri {
				t.Errorf("Expected URI %s, got %s", tt.uri, doc.URI)
			}

			if doc.Basepath != tt.wantPath {
				t.Errorf("Expected basepath %s, got %s", tt.wantPath, doc.Basepath)
			}
		})
	}
}
