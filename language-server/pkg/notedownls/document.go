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
	// Content is the current text content of the document
	Content string
	// Version is the current version number of the document
	Version int
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
		Content:  "",
		Version:  0,
	}, nil
}

// NewDocumentWithContent creates a new Document with initial content
func NewDocumentWithContent(uri, content string, version int) (*Document, error) {
	basepath, err := uriToBasepath(uri)
	if err != nil {
		return nil, err
	}

	return &Document{
		URI:      uri,
		Basepath: basepath,
		Content:  content,
		Version:  version,
	}, nil
}

// UpdateContent updates the document content and version
func (d *Document) UpdateContent(content string, version int) {
	d.Content = content
	d.Version = version
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
