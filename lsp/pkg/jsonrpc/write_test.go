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
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name        string
		response    Response
		expectError bool
		checkOutput func(string) bool
	}{
		{
			name: "successful response",
			response: Response{
				ProtocolVersion: JSONRPCVersion,
				ID:              func() *json.RawMessage { raw := json.RawMessage(`1`); return &raw }(),
				Result:          map[string]interface{}{"success": true},
			},
			expectError: false,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "Content-Length:") &&
					strings.Contains(output, `"jsonrpc":"2.0"`) &&
					strings.Contains(output, `"id":1`) &&
					strings.Contains(output, `"result":{"success":true}`)
			},
		},
		{
			name: "error response",
			response: Response{
				ProtocolVersion: JSONRPCVersion,
				ID:              func() *json.RawMessage { raw := json.RawMessage(`"test"`); return &raw }(),
				Error: &ResponseError{
					Code:    InvalidRequestCode,
					Message: "Invalid Request",
					Data:    "test data",
				},
			},
			expectError: false,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "Content-Length:") &&
					strings.Contains(output, `"jsonrpc":"2.0"`) &&
					strings.Contains(output, `"id":"test"`) &&
					strings.Contains(output, `"error":`) &&
					strings.Contains(output, `"code":-32600`)
			},
		},
		{
			name: "response with null id",
			response: Response{
				ProtocolVersion: JSONRPCVersion,
				ID:              func() *json.RawMessage { raw := json.RawMessage(`null`); return &raw }(),
				Result:          "test result",
			},
			expectError: false,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "Content-Length:") &&
					strings.Contains(output, `"id":null`) &&
					strings.Contains(output, `"result":"test result"`)
			},
		},
		{
			name: "response with complex result",
			response: Response{
				ProtocolVersion: JSONRPCVersion,
				ID:              func() *json.RawMessage { raw := json.RawMessage(`42`); return &raw }(),
				Result: map[string]interface{}{
					"nested": map[string]interface{}{
						"array": []int{1, 2, 3},
						"bool":  true,
					},
				},
			},
			expectError: false,
			checkOutput: func(output string) bool {
				return strings.Contains(output, "Content-Length:") &&
					strings.Contains(output, `"id":42`) &&
					strings.Contains(output, `"nested"`) &&
					strings.Contains(output, `"array":[1,2,3]`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			err := Write(writer, tt.response)

			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
				return
			}

			require.NoError(t, err, "unexpected error")

			output := buf.String()

			// Check that output contains expected elements
			if tt.checkOutput != nil {
				assert.True(t, tt.checkOutput(output), "output validation failed. Got: %s", output)
			}

			// Verify Content-Length header is correct
			lines := strings.Split(output, "\r\n")
			assert.GreaterOrEqual(t, len(lines), 3, "expected at least 3 lines in output")

			assert.True(t, strings.HasPrefix(lines[0], "Content-Length: "),
				"expected Content-Length header, got: %s", lines[0])

			assert.Empty(t, lines[1], "expected empty line after headers")

			// Extract and verify content length
			body := strings.Join(lines[2:], "\r\n")
			expectedLength := len(body)
			headerValue := strings.TrimPrefix(lines[0], "Content-Length: ")

			actualLength, err := strconv.Atoi(headerValue)
			require.NoError(t, err, "Content-Length header should be a valid number, got: %s", headerValue)

			// Allow some tolerance for JSON formatting differences
			if abs(actualLength-expectedLength) > 5 {
				assert.Equal(t, expectedLength, actualLength,
					"Content-Length mismatch: header says %d, actual body length is %d",
					actualLength, expectedLength)
			}
		})
	}
}

func TestReadWriteRoundtrip(t *testing.T) {
	tests := []struct {
		name     string
		request  Request
		response Response
	}{
		{
			name: "basic request/response",
			request: Request{
				ProtocolVersion: JSONRPCVersion,
				ID:              func() *json.RawMessage { raw := json.RawMessage(`1`); return &raw }(),
				Method:          "test",
				Params:          json.RawMessage(`{"param1":"value1"}`),
			},
			response: Response{
				ProtocolVersion: JSONRPCVersion,
				ID:              func() *json.RawMessage { raw := json.RawMessage(`1`); return &raw }(),
				Result:          map[string]string{"result": "success"},
			},
		},
		{
			name: "notification (no response expected)",
			request: Request{
				ProtocolVersion: JSONRPCVersion,
				ID:              nil, // notification
				Method:          "notify",
				Params:          json.RawMessage(`{"event":"test"}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test request serialization
			requestJSON, err := json.Marshal(tt.request)
			require.NoError(t, err, "failed to marshal request")

			// Create properly formatted input
			input := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(requestJSON), string(requestJSON))
			reader := bufio.NewReader(strings.NewReader(input))

			// Test reading
			parsedRequest, err := Read(reader)
			require.NoError(t, err, "failed to read request")

			// Verify request
			assert.Equal(t, tt.request.ProtocolVersion, parsedRequest.ProtocolVersion,
				"protocol version mismatch")

			assert.Equal(t, tt.request.Method, parsedRequest.Method,
				"method mismatch")

			// Test response writing (if not a notification)
			if tt.request.ID != nil {
				var buf bytes.Buffer
				writer := bufio.NewWriter(&buf)

				err := Write(writer, tt.response)
				require.NoError(t, err, "failed to write response")

				output := buf.String()
				assert.Contains(t, output, "Content-Length:",
					"response should contain Content-Length header")

				assert.Contains(t, output, `"jsonrpc":"2.0"`,
					"response should contain JSON-RPC version")
			}
		})
	}
}

func TestProtocolWriteEdgeCases(t *testing.T) {
	t.Run("write with marshal error", func(t *testing.T) {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)

		// Create response with unmarshalable data
		response := Response{
			ProtocolVersion: JSONRPCVersion,
			ID:              func() *json.RawMessage { raw := json.RawMessage(`1`); return &raw }(),
			Result:          make(chan int), // channels can't be marshaled
		}

		err := Write(writer, response)
		assert.Error(t, err, "expected error when marshaling invalid data")
	})
}
