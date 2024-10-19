package tasks_test

import (
	"fmt"

	"github.com/liamawhite/nl/pkg/workspace/documents/writer"
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
