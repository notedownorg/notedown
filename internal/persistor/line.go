package persistor

import (
	"fmt"
	"os"
	"strings"
)

// AddLine adds a line of text to a document at the specified line number.
// If the line number is -1 or greater than the number of lines in the document,
// the line is added to the end of the document.
func (p *Persistor) AddLine(document string, line int, text string) error {
	err := p.validateLine(text)
	if err != nil {
		return fmt.Errorf("invalid text: %w", err)
	}

	lines, err := p.readLines(document)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	if line == -1 || line >= len(lines) {
		lines = append(lines, text)
	} else {
		lines = append(lines[:line], append([]string{text}, lines[line:]...)...)
	}

	err = p.writeLines(document, lines)
	if err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

// RemoveLine removes a line of text from a document at the specified line number.
func (p *Persistor) RemoveLine(document string, line int) error {
	lines, err := p.readLines(document)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	if line < 0 || line >= len(lines) {
		return fmt.Errorf("line number out of bounds")
	}

	lines = append(lines[:line], lines[line+1:]...)

	err = p.writeLines(document, lines)
	if err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

// UpdateLine updates a line of text in a document at the specified line number.
func (p *Persistor) UpdateLine(document string, line int, text string) error {
	err := p.validateLine(text)
	if err != nil {
		return fmt.Errorf("invalid text: %w", err)
	}

	lines, err := p.readLines(document)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	if line < 0 || line >= len(lines) {
		return fmt.Errorf("line number out of bounds")
	}

	lines[line] = text

	err = p.writeLines(document, lines)
	if err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

// opens a document and returns its lines.
func (p *Persistor) readLines(document string) ([]string, error) {
	bytes, err := os.ReadFile(document)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), "\n"), nil
}

func (p *Persistor) writeLines(document string, lines []string) error {
	return os.WriteFile(document, []byte(strings.Join(lines, "\n")), 0644)
}

func (p *Persistor) validateLine(text string) error {
	if strings.Contains(text, "\n") {
		return fmt.Errorf("text contains newline character")
	}
	return nil
}
