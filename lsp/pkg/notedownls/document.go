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
	// Basepath is the filename without extension (e.g., "README" for "README.md")
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

// uriToBasepath converts a file:// URI to filename without extension
func uriToBasepath(uri string) (string, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if parsedURI.Scheme != "file" {
		return "", nil // Non-file URIs don't have basepaths
	}

	// Get the filename from the path
	filename := filepath.Base(parsedURI.Path)

	// Remove the extension to get the base name
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))

	return baseName, nil
}
