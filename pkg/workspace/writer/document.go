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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/notedownorg/notedown/pkg/workspace"
)

type FileExistsError struct {
	Filename string
}

func (e *FileExistsError) Error() string {
	return fmt.Sprintf("file %s already exists", e.Filename)
}

func (c Client) Create(doc workspace.Document) error {
	slog.Debug("creating document", "path", doc.Path())

	// Ensure the document is in a valid state
	workspace.Sanitize(*c.config, &doc)

	// Ensure the file does not exist
	_, err := os.Stat(c.abs(doc.Path()))
	if err == nil {
		return &FileExistsError{Filename: doc.Path()}
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(c.abs(doc.Path())), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	return os.WriteFile(c.abs(doc.Path()), []byte(doc.Markdown()), 0644)
}

func (c Client) Update(doc workspace.Document) error {
	slog.Debug("updating document", "path", doc.Path())

	// Ensure the document is in a valid state
	workspace.Sanitize(*c.config, &doc)

	// Ensure the file exists
	stat, err := os.Stat(c.abs(doc.Path()))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", doc.Path())
		}
	}

	// Ensure the file hasnt changed since we loaded it into memory
	// This is a simple check to ensure we're not accidentally overwriting changes made by another process
	if doc.Modified(stat.ModTime()) {
		return fmt.Errorf("file %s has been modified since it was loaded", doc.Path())
	}

	// Write the file
	return os.WriteFile(c.abs(doc.Path()), []byte(doc.Markdown()), 0644)
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

func (c Client) Delete(path string) error {
	slog.Debug("deleting document", "path", path)
	if err := os.Remove(c.abs(path)); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete document: %w", err)
		}
	}
	return nil
}
