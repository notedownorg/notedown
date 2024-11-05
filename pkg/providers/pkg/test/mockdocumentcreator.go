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

type AddValidator func(doc writer.Document, metadata reader.Metadata, content []byte, feed chan reader.Event) error

type MockDocumentCreator struct {
	Feed       chan reader.Event
	Validators []AddValidator
}

func (m *MockDocumentCreator) Add(path string, metadata reader.Metadata, content []byte) error {
	return m.validate(writer.Document{Path: path}, metadata, content)
}

func (m *MockDocumentCreator) validate(doc writer.Document, metadata reader.Metadata, content []byte) error {
	if len(m.Validators) == 0 {
		return fmt.Errorf("no validators left")
	}
	validator := m.Validators[0]
	m.Validators = m.Validators[1:]
	slog.Info("removed validator", "remaining", len(m.Validators), "doc", doc, "metadata", metadata, "content", content)
	return validator(doc, metadata, content, m.Feed)
}
