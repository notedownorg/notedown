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

package workspace

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	"github.com/notedownorg/notedown/pkg/parse/blocks"
)

const (
	MetadataTagsKey = "tags"
)

type Metadata map[string]interface{}

func (m Metadata) Types() []string {
	tags := make([]string, 0)
	if m == nil {
		return tags
	}
	tagsValue, ok := m[MetadataTagsKey]
	if !ok {
		return tags
	}
	res, ok := tagsValue.([]string)
	if !ok {
		// Check if tags is set to a single string instead of a list
		tagsValue, ok = m[MetadataTagsKey].(string)
		if !ok {
			return tags
		}
		res = []string{tagsValue.(string)}
	}

	for _, tag := range res {
		split := strings.Split(tag, "/")
		tags = append(tags, split[0])
	}
	return tags
}

func (m Metadata) HasType(t string) bool {
	types := m.Types()
	for _, typ := range types {
		if typ == t {
			return true
		}
	}
	return false
}

func NewMetadata() Metadata {
	return Metadata{MetadataTagsKey: make([]string, 0)}
}

type Document struct {
	// Moving Metadata to the top of the struct is useful for removing the amount of defensive code by consumers
	Metadata Metadata
	Blocks   []ast.Block

	path string

	// Do not modify these fields after creation, used to detect mutations on disk and in memory after creation.
	// We require the hash as well as the file mod time because we don't have identical round-tripping and need to tolerate that rather than explode.
	lastModified time.Time
	creationHash string
}

// Use when creating a new document
func NewDocument(relativePath string, metadata Metadata, blocks ...ast.Block) Document {
	return Document{path: relativePath, Metadata: metadata, Blocks: blocks}
}

// Load and parse a document from a file
func LoadDocument(workspaceDir string, relativePath string, relativeTo time.Time) (Document, error) {
	fullPath := filepath.Join(workspaceDir, relativePath)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return Document{}, fmt.Errorf("failed to get file info: %w", err)
	}
	input, err := os.ReadFile(fullPath)
	if err != nil {
		return Document{}, fmt.Errorf("failed to read file: %w", err)
	}

	p := NewInput(string(input))
	res, ok, err := documentParser(relativeTo)(p)
	if err != nil {
		return Document{}, fmt.Errorf("unable to parse document: %w", err)
	}
	if !ok {
		return Document{}, fmt.Errorf("unable to parse document")
	}

	res.lastModified = stat.ModTime()
	res.creationHash = hash(res.Markdown())
	res.path = relativePath
	return res, nil
}

func (d Document) Markdown() string {
	var output strings.Builder
	if d.Metadata != nil {
		metadata := blocks.NewFrontmatter(d.Metadata)
		output.WriteString(metadata.Markdown())
		output.WriteString("\n")
	}
	for _, block := range d.Blocks {
		output.WriteString(block.Markdown())
		output.WriteString("\n")
	}
	return output.String()
}

func (d Document) Path() string {
	return d.path
}

func hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Has the document been changed in memory (since read time)
func (d Document) Mutated() bool {
	return d.creationHash != hash(d.Markdown())
}

// Is the document in memory stale compared to the file on disk
func (d Document) Modified(lastModified time.Time) bool {
	if d.path == "basic.md" {
		fmt.Println("lastModified: ", lastModified.UnixNano())
		fmt.Println("d.lastModified: ", d.lastModified.UnixNano())
	}
	return lastModified.UnixNano() > d.lastModified.UnixNano()
}
