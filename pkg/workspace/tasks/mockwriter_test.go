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

package tasks_test

import (
	"fmt"

	"github.com/notedownorg/notedown/pkg/workspace/documents/writer"
)

var _ writer.LineWriter = &MockLineWriter{}

type validator func(method string, doc writer.Document, line int, obj fmt.Stringer) error

type MockLineWriter struct {
	validators []validator
}

func (m *MockLineWriter) validate(method string, doc writer.Document, line int, obj fmt.Stringer) error {
	if len(m.validators) == 0 {
		return fmt.Errorf("no validators left")
	}
	validator := m.validators[0]
	m.validators = m.validators[1:]
	return validator(method, doc, line, obj)
}

func (m *MockLineWriter) AddLine(doc writer.Document, line int, obj fmt.Stringer) error {
	return m.validate("add", doc, line, obj)
}

func (m *MockLineWriter) RemoveLine(doc writer.Document, line int) error {
	return m.validate("remove", doc, line, nil)
}

func (m *MockLineWriter) UpdateLine(doc writer.Document, line int, obj fmt.Stringer) error {
	return m.validate("update", doc, line, obj)
}
