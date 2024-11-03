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

func (c Client) AddDocument(path string, metadata reader.Metadata, content []byte) error {
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

// func (c Client) RemoveDocument(doc Document) error {
//     return nil
// }
//
// func (c Client) UpdateDocument(doc Document, metadata reader.Metadata, content []byte) error {
//     return nil
// }
