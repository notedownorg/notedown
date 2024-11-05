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
	"fmt"
	"math"
	"strings"
)

const (
	AT_BEGINNING int = 0
	AT_END       int = math.MaxInt
)

type LineMutation func(checksum string, lines []string) ([]string, error)

// Add a line to the content. All line numbers are 1-indexed and do not include frontmatter
// e.g. if the frontmatter ends at line 3 in the underlying file, adding a line at 1 will insert the line at 4
func AddLine(number int, obj fmt.Stringer) LineMutation {
	return func(checksum string, lines []string) ([]string, error) {
		// Validate that the inputs
		if err := validateLine(obj.String()); err != nil {
			return lines, fmt.Errorf("invalid text adding line '%s': %w", obj, err)
		}
		if checksum == "" && number != AT_END && number != AT_BEGINNING {
			return lines, fmt.Errorf("hash must be provided when adding a line in the middle of a document")
		}

		// If we're adding at the end, just append the line
		if number == AT_END || number > len(lines) {
			return append(lines, obj.String()), nil
		}

		// Handle adding at the beginning
		if number == AT_BEGINNING || number <= 1 {
			return append([]string{obj.String()}, lines...), nil
		}

		// 0-indexed but input is 1-indexed
		return append(lines[:number-1], append([]string{obj.String()}, lines[number-1:]...)...), nil
	}
}

// Remove a line from the content. All line numbers are 1-indexed and do not include frontmatter
// e.g. if the frontmatter ends at line 3 in the underlying file, removing a line at 1 will remove the line at 4
func RemoveLine(number int) LineMutation {
	return func(checksum string, lines []string) ([]string, error) {
		if checksum == "" {
			return lines, fmt.Errorf("hash must be provided when removing a line to avoid stale writes")
		}
		if number == AT_END || number == AT_BEGINNING {
			return lines, fmt.Errorf("must provide an absolute line number when removing a line")
		}

		if number <= 0 || number > len(lines) {
			return lines, fmt.Errorf("line number out of bounds")
		}

		// 0-indexed but input is 1-indexed
		return append(lines[:number-1], lines[number:]...), nil
	}
}

// Update a line in the content. All line numbers are 1-indexed and do not include frontmatter
// e.g. if the frontmatter ends at line 3 in the underlying file, updating a line at 1 will update the line at 4
func UpdateLine(number int, obj fmt.Stringer) LineMutation {
	return func(checksum string, lines []string) ([]string, error) {
		if err := validateLine(obj.String()); err != nil {
			return lines, fmt.Errorf("invalid text updating line '%s': %w", obj, err)
		}
		if checksum == "" {
			return lines, fmt.Errorf("hash must be provided when updating a line to avoid stale writes")
		}
		if number == AT_END || number == AT_BEGINNING {
			return lines, fmt.Errorf("must provide an absolute line number")
		}

		if number <= 0 || number > len(lines) {
			return lines, fmt.Errorf("line number out of bounds")
		}

		// 0-indexed but input is 1-indexed
		lines[number-1] = obj.String()
		return lines, nil
	}
}

func validateLine(text string) error {
	if strings.Contains(text, "\n") {
		return fmt.Errorf("text contains newline character")
	}
	return nil
}
