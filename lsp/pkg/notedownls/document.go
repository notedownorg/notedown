package notedownls

import (
	"net/url"
	"path/filepath"
	"strings"
)

// Document represents a text document being tracked by the language server
type Document struct {
	// URI is the document identifier as per LSP specification
	URI string
	// Basepath is the file path relative to the workspace root, derived from URI
	Basepath string
}

// NewDocument creates a new Document from a URI
func NewDocument(uri string) (*Document, error) {
	basepath, err := uriToBasepath(uri)
	if err != nil {
		return nil, err
	}

	return &Document{
		URI:      uri,
		Basepath: basepath,
	}, nil
}

// uriToBasepath converts a file:// URI to a relative basepath
func uriToBasepath(uri string) (string, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if parsedURI.Scheme != "file" {
		return "", nil // Non-file URIs don't have basepaths
	}
	path := strings.TrimPrefix(parsedURI.Path, "/")
	return filepath.Clean(path), nil
}
