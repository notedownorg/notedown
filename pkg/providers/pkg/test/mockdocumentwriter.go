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

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
)

var _ writer.DocumentWriter = &MockDocumentWriter{}

type DocumentWriterValidator func(method string, doc writer.Document, metadata reader.Metadata, content []byte, feed chan reader.Event) error

type MockDocumentWriter struct {
	Feed       chan reader.Event
	Validators []DocumentWriterValidator
}

func (m *MockDocumentWriter) AddDocument(path string, metadata reader.Metadata, content []byte) error {
	return m.validate("add", writer.Document{Path: path}, metadata, content)
}

func (m *MockDocumentWriter) validate(method string, doc writer.Document, metadata reader.Metadata, content []byte) error {
	if len(m.Validators) == 0 {
		return fmt.Errorf("no validators left")
	}
	validator := m.Validators[0]
	m.Validators = m.Validators[1:]
	return validator(method, doc, metadata, content, m.Feed)
}
