package writer

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

const (
	AtBeginning int = 0
	AtEnd       int = -1
)

// AddLine adds a line of text to a document at the specified line number.
func (c Client) AddLine(doc Document, line int, obj fmt.Stringer) error {
	err := validateLine(obj.String())
	if err != nil {
		return fmt.Errorf("invalid text: %w", err)
	}

	lines, frontmatter, err := readAndValidateFile(c.abs(doc.Path), doc.Hash)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	// If we're within the frontmatter return an error
	// Unless we're adding AtBeginning, which should add the line after the frontmatter
	if frontmatter != -1 && line <= frontmatter && line != AtBeginning {
		return fmt.Errorf("cannot add a line within frontmatter")
	}

	if line == AtEnd || line > len(lines) {
		lines = append(lines, obj.String())
	} else if line == AtBeginning {
		if frontmatter != -1 {
			// If the file has frontmatter, we need to insert the line after it
			lines = append(lines[:frontmatter], append([]string{obj.String()}, lines[frontmatter:]...)...)
		} else {
			lines = append([]string{obj.String()}, lines...)
		}
	} else {
		// 0-indexed but input is 1-indexed
		lines = append(lines[:line-1], append([]string{obj.String()}, lines[line:]...)...)
	}

	err = writeLines(c.abs(doc.Path), lines)
	if err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

// RemoveLine removes a line of text from a document at the specified line number.
func (c Client) RemoveLine(doc Document, line int) error {
	lines, frontmatter, err := readAndValidateFile(c.abs(doc.Path), doc.Hash)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	if line == AtEnd || line == AtBeginning {
		return fmt.Errorf("must provide an absolute line number")
	}

	if line <= 0 || line > len(lines) {
		return fmt.Errorf("line number out of bounds")
	}

	// If we're within the frontmatter return an error
	if frontmatter != -1 && line <= frontmatter {
		return fmt.Errorf("cannot remove a line within frontmatter")
	}

	// 0-indexed but input is 1-indexed
	lines = append(lines[:line-1], lines[line:]...)

	err = writeLines(c.abs(doc.Path), lines)
	if err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

// UpdateLine updates a line of text in a document at the specified line number.
func (c Client) UpdateLine(doc Document, line int, obj fmt.Stringer) error {
	err := validateLine(obj.String())
	if err != nil {
		return fmt.Errorf("invalid text: %w", err)
	}

	lines, frontmatter, err := readAndValidateFile(c.abs(doc.Path), doc.Hash)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	if line == AtEnd || line == AtBeginning {
		return fmt.Errorf("must provide an absolute line number")
	}

	if line <= 0 || line > len(lines) {
		return fmt.Errorf("line number out of bounds")
	}

	// If we're within the frontmatter return an error
	if frontmatter != -1 && line <= frontmatter {
		return fmt.Errorf("cannot update a line within frontmatter")
	}

	lines[line-1] = obj.String() // 0-indexed

	err = writeLines(c.abs(doc.Path), lines)
	if err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

func writeLines(path string, lines []string) error {
	// maintain trailing newline
	content := bytes.NewBuffer([]byte{})
	for _, line := range lines {
		content.WriteString(line)
		content.WriteString("\n")
	}

	return os.WriteFile(path, content.Bytes(), 0644)
}

func validateLine(text string) error {
	if strings.Contains(text, "\n") {
		return fmt.Errorf("text contains newline character")
	}
	return nil
}

// returns lines, where the frontmatter ends or -1 if there is no frontmatter and an error
func readAndValidateFile(path string, hash string) ([]string, int, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, -1, err
	}

	algo := sha256.New()
	algo.Write(bytes)
	if hash != fmt.Sprintf("%x", algo.Sum(nil)) {
		return nil, -1, fmt.Errorf("file has been modified since last read, unable to write with stale data")
	}

	lines := strings.Split(string(bytes), "\n")

	// Remove the last line if it's empty to prevent adding additional whitespace
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Check if the file has frontmatter
	// This is a simple/fast check, but it should be sufficient for this use case
	frontmatter := -1
	if len(lines) > 0 && strings.HasPrefix(lines[0], "---") {
		for i, line := range lines[1:] {
			if strings.HasPrefix(line, "---") {
				frontmatter = i + 2 // 0 -> 1-indexed and after the current line
				break
			}
		}
	}

	return lines, frontmatter, nil
}
