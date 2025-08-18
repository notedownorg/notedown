// Copyright 2025 Notedown Authors
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

package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
)

func Write(w *bufio.Writer, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}
	headers := fmt.Sprintf("%s: %d\r\n\r\n", ContentLengthHeader, len(body))
	if _, err := w.WriteString(headers); err != nil {
		return fmt.Errorf("error writing headers: %w", err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("error writing body: %w", err)
	}
	return w.Flush()
}
