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

	"github.com/notedownorg/notedown/pkg/workspace"
	"github.com/notedownorg/notedown/pkg/workspace/reader"
)

type CreateValidator func(doc workspace.Document) error
type RenameValidator func(oldPath, newPath string) error
type DeleteValidator func(path string) error

type MockDocumentWriter struct {
	Feed       chan reader.Event
	Validators Validators
}

type Validators struct {
	Create []CreateValidator
	Rename []RenameValidator
	Delete []DeleteValidator
}

func (m *MockDocumentWriter) Create(doc workspace.Document) error {
	return m.validateCreate(doc)
}

func (m *MockDocumentWriter) validateCreate(doc workspace.Document) error {
	if len(m.Validators.Create) == 0 {
		return fmt.Errorf("no create validators left")
	}
	validator := m.Validators.Create[0]
	m.Validators.Create = m.Validators.Create[1:]
	slog.Info("removed create validator", "remaining", len(m.Validators.Create), "doc", doc.Path())
	return validator(doc)
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

func (m *MockDocumentWriter) Delete(path string) error {
	return m.validateDelete(path)
}

func (m *MockDocumentWriter) validateDelete(path string) error {
	if len(m.Validators.Delete) == 0 {
		return fmt.Errorf("no delete validators left")
	}
	validator := m.Validators.Delete[0]
	m.Validators.Delete = m.Validators.Delete[1:]
	return validator(path)
}
