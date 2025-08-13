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
			wantPath: "path/to/document.md",
		},
		{
			name:     "file URI with relative path",
			uri:      "file://./README.md",
			wantPath: "README.md",
		},
		{
			name:     "non-file URI",
			uri:      "http://example.com/doc",
			wantPath: "",
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
