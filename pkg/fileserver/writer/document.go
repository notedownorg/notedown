// Copyright 2024 Notedown Authors
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

package writer

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"sigs.k8s.io/yaml"
)

type Document struct {
	// Path is relative to root
	Path     string
	Checksum string
}

type FileExistsError struct {
	Filename string
}

func (e *FileExistsError) Error() string {
	return fmt.Sprintf("file %s already exists", e.Filename)
}

func (c Client) Create(path string, metadata reader.Metadata, content []byte) error {
	// Ensure the file does not exist
	_, err := os.Stat(c.abs(path))
	if err == nil {
		return &FileExistsError{Filename: path}
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	var b bytes.Buffer
	if metadata != nil && len(metadata) > 0 {
		md, err := yaml.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		b.WriteString("---\n")
		b.Write(md)
		b.WriteString("---\n")
	}
	b.Write(content)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(c.abs(path)), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	return os.WriteFile(c.abs(path), b.Bytes(), 0644)
}

func (c Client) UpdateMetadata(doc Document, metadata reader.Metadata) error {
	slog.Debug("updating metadata of document", "path", doc.Path)

	lines, frontmatter, err := readAndValidateFile(c.abs(doc.Path), doc.Checksum)
	if err != nil {
		return fmt.Errorf("failed to validate document: %w", err)
	}

	// Drop the current frontmatter if it exists
	if frontmatter != -1 {
		lines = lines[frontmatter:]
	}

	// Update the metadata
	var b bytes.Buffer
	if metadata != nil && len(metadata) > 0 {
		md, err := yaml.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		b.WriteString("---\n")
		b.Write(md)
		b.WriteString("---\n")
	}

	content := bytes.NewBuffer([]byte{})
	content.Write(b.Bytes())
	for _, line := range lines {
		content.WriteString(line)
		content.WriteString("\n")
	}

	if err := os.WriteFile(c.abs(doc.Path), content.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}

	return nil
}

// Update contents of a document. Mutations are applied in order and are atomeic.
// i.e. if any mutation errors, the document will not be written to disk.
func (c Client) UpdateContent(doc Document, mutations ...LineMutation) error {
	slog.Debug("updating content of document", "path", doc.Path)

	lines, frontmatter, err := readAndValidateFile(c.abs(doc.Path), doc.Checksum)
	if err != nil {
		return fmt.Errorf("failed to validate document: %w", err)
	}

	// Split out the frontmatter if it exists as mutations assume it's not there
	prefix := make([]string, 0)
	if frontmatter != -1 {
		prefix, lines = lines[:frontmatter], lines[frontmatter:]
	}

	for i, mutation := range mutations {
		lines, err = mutation(doc.Checksum, lines)
		if err != nil {
			return fmt.Errorf("invalid line mutation at index %d, no mutations will be written to disk: %w", i, err)
		}
	}

	content := bytes.NewBuffer([]byte{})
	for _, line := range append(prefix, lines...) {
		content.WriteString(line)
		content.WriteString("\n")
	}

	if err := os.WriteFile(c.abs(doc.Path), content.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}

	return nil
}

func (c Client) Rename(oldPath, newPath string) error {
	slog.Debug("renaming document", "oldPath", oldPath, "newPath", newPath)

	// Check if the new path already exists
	_, err := os.Stat(c.abs(newPath))
	if err == nil {
		return &FileExistsError{Filename: newPath}
	}

	// As we're not changing the content we don't need to validate the checksum
	if err := os.Rename(c.abs(oldPath), c.abs(newPath)); err != nil {
		return fmt.Errorf("failed to rename document: %w", err)
	}
	return nil
}

func (c Client) Delete(doc Document) error {
	slog.Debug("deleting document", "path", doc.Path)
	if err := os.Remove(c.abs(doc.Path)); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete document: %w", err)
		}
	}
	return nil
}
