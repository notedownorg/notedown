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
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

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
		latest := fmt.Sprintf("%x", algo.Sum(nil))
		if checksum != latest {
			return nil, -1, fmt.Errorf("file has been modified since last read, unable to write with stale data wanted: %s got: %s", latest, checksum)
		}
	}

	res := strings.Split(string(bytes), "\n")

	// Remove the last line if it's empty to prevent adding additional whitespace
	if len(res) > 0 && res[len(res)-1] == "" {
		res = res[:len(res)-1]
	}

	// Check if the file has frontmatter
	// This is a simple/fast check, but it should be sufficient for this use case
	frontmatter := -1
	if len(res) > 0 && strings.HasPrefix(res[0], "---") {
		for i, line := range res[1:] {
			if strings.HasPrefix(line, "---") {
				frontmatter = i + 2 // 0 -> 1-indexed and after the current line
				break
			}
		}
	}

	return res, frontmatter, nil
}
