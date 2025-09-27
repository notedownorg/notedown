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

package notedownls

import (
	"encoding/json"
	"fmt"

	"github.com/notedownorg/notedown/language-server/pkg/lsp"
)

// handleExecuteCommand handles workspace/executeCommand requests
func (s *Server) handleExecuteCommand(params json.RawMessage) (any, error) {
	var executeParams lsp.ExecuteCommandParams
	if err := json.Unmarshal(params, &executeParams); err != nil {
		s.logger.Error("failed to unmarshal execute command params", "error", err)
		return nil, err
	}

	s.logger.Debug("execute command request received", "command", executeParams.Command)

	switch executeParams.Command {
	case "notedown.getListItemBoundaries":
		return s.handleGetListItemBoundaries(executeParams.Arguments)
	case "notedown.getConcealRanges":
		return s.handleGetConcealRanges(executeParams.Arguments)
	case "notedown.executeCodeBlocks":
		// Convert []any to []json.RawMessage
		var rawArgs []json.RawMessage
		for _, arg := range executeParams.Arguments {
			if argBytes, err := json.Marshal(arg); err == nil {
				rawArgs = append(rawArgs, argBytes)
			}
		}
		return s.handleExecuteCodeBlocks(rawArgs)
	default:
		return nil, fmt.Errorf("unknown command: %s", executeParams.Command)
	}
}
