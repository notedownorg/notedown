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
	"crypto/sha256"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strings"
)

const (
	AT_BEGINNING int = 0
	AT_END       int = math.MaxInt
)

// AddLine adds a line of text to a document at the specified line number.
func (c Client) AddLine(doc Document, line int, obj fmt.Stringer) error {
	slog.Debug("adding line to document", "path", doc.Path, "number", line, "text", obj)
	if doc.Checksum == "" && line != AT_END && line != AT_BEGINNING {
		return fmt.Errorf("hash must be provided when adding a line in the middle of a document")
	}

	err := validateLine(obj.String())
	if err != nil {
		return fmt.Errorf("invalid text: %w", err)
	}

	lines, frontmatter, err := readAndValidateFile(c.abs(doc.Path), doc.Checksum)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	// If we're within the frontmatter return an error
	// Unless we're adding AtBeginning, which should add the line after the frontmatter
	if frontmatter != -1 && line <= frontmatter && line != AT_BEGINNING {
		return fmt.Errorf("cannot add a line within frontmatter")
	}

	if line == AT_END || line > len(lines) {
		lines = append(lines, obj.String())
	} else if line == AT_BEGINNING {
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
	slog.Debug("removing line from document", "path", doc.Path, "number", line)
	if doc.Checksum == "" {
		return fmt.Errorf("hash must be provided when removing a line to avoid stale writes")
	}

	lines, frontmatter, err := readAndValidateFile(c.abs(doc.Path), doc.Checksum)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	if line == AT_END || line == AT_BEGINNING {
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
	slog.Debug("updating line in document", "path", doc.Path, "number", line, "text", obj)
	if doc.Checksum == "" {
		return fmt.Errorf("hash must be provided when updating a line to avoid stale writes")
	}

	err := validateLine(obj.String())
	if err != nil {
		return fmt.Errorf("invalid text: %w", err)
	}

	lines, frontmatter, err := readAndValidateFile(c.abs(doc.Path), doc.Checksum)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}

	if line == AT_END || line == AT_BEGINNING {
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
func readAndValidateFile(path string, checksum string) ([]string, int, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, -1, err
	}

	// Ensure file hasn't been modified only if a hash is provided
	if checksum != "" {
		algo := sha256.New()
		algo.Write(bytes)
		if checksum != fmt.Sprintf("%x", algo.Sum(nil)) {
			return nil, -1, fmt.Errorf("file has been modified since last read, unable to write with stale data")
		}
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
