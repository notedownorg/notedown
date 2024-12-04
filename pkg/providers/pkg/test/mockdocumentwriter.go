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

package test

import (
	"fmt"
	"log/slog"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
)

type CreateValidator func(doc writer.Document, metadata reader.Metadata, content []byte, feed chan reader.Event) error
type MetadataUpdateValidator func(doc writer.Document, metadata reader.Metadata) error
type ContentUpdateValidator func(doc writer.Document, mutations ...writer.LineMutation) error
type RenameValidator func(oldPath, newPath string) error
type DeleteValidator func(doc writer.Document) error

type MockDocumentWriter struct {
	Feed       chan reader.Event
	Validators Validators
}

type Validators struct {
	Create         []CreateValidator
	MetadataUpdate []MetadataUpdateValidator
	ContentUpdate  []ContentUpdateValidator
	Rename         []RenameValidator
	Delete         []DeleteValidator
}

func (m *MockDocumentWriter) Create(path string, metadata reader.Metadata, content []byte) error {
	return m.validateCreate(writer.Document{Path: path}, metadata, content)
}

func (m *MockDocumentWriter) validateCreate(doc writer.Document, metadata reader.Metadata, content []byte) error {
	if len(m.Validators.Create) == 0 {
		return fmt.Errorf("no create validators left")
	}
	validator := m.Validators.Create[0]
	m.Validators.Create = m.Validators.Create[1:]
	slog.Info("removed create validator", "remaining", len(m.Validators.Create), "doc", doc, "metadata", metadata, "content", content)
	return validator(doc, metadata, content, m.Feed)
}

func (m *MockDocumentWriter) UpdateMetadata(doc writer.Document, metadata reader.Metadata) error {
	return m.validateMetadataUpdate(doc, metadata)
}

func (m *MockDocumentWriter) validateMetadataUpdate(doc writer.Document, metadata reader.Metadata) error {
	if len(m.Validators.MetadataUpdate) == 0 {
		return fmt.Errorf("no metadata update validators left")
	}
	validator := m.Validators.MetadataUpdate[0]
	m.Validators.MetadataUpdate = m.Validators.MetadataUpdate[1:]
	return validator(doc, metadata)
}

func (m *MockDocumentWriter) UpdateContent(doc writer.Document, mutations ...writer.LineMutation) error {
	return m.validateContentUpdate(doc, mutations...)
}

func (m *MockDocumentWriter) validateContentUpdate(doc writer.Document, mutations ...writer.LineMutation) error {
	if len(m.Validators.ContentUpdate) == 0 {
		return fmt.Errorf("no content update validators left")
	}
	validator := m.Validators.ContentUpdate[0]
	m.Validators.ContentUpdate = m.Validators.ContentUpdate[1:]
	return validator(doc, mutations...)
}

func (m *MockDocumentWriter) Rename(oldPath, newPath string) error {
	return m.validateRename(oldPath, newPath)
}

func (m *MockDocumentWriter) validateRename(oldPath, newPath string) error {
	if len(m.Validators.Rename) == 0 {
		return fmt.Errorf("no rename validators left")
	}
	validator := m.Validators.Rename[0]
	m.Validators.Rename = m.Validators.Rename[1:]
	return validator(oldPath, newPath)
}

func (m *MockDocumentWriter) Delete(doc writer.Document) error {
	return m.validateDelete(doc)
}

func (m *MockDocumentWriter) validateDelete(doc writer.Document) error {
	if len(m.Validators.Delete) == 0 {
		return fmt.Errorf("no delete validators left")
	}
	validator := m.Validators.Delete[0]
	m.Validators.Delete = m.Validators.Delete[1:]
	return validator(doc)
}
