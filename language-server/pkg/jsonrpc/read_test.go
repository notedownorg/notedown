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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectError  bool
		expectID     *json.RawMessage
		expectMethod string
	}{
		{
			name: "valid request with id",
			input: "Content-Length: 67\r\n\r\n" +
				`{"jsonrpc":"2.0","id":1,"method":"test","params":{"key":"value"}}`,
			expectError:  false,
			expectID:     func() *json.RawMessage { raw := json.RawMessage(`1`); return &raw }(),
			expectMethod: "test",
		},
		{
			name: "valid notification without id",
			input: "Content-Length: 61\r\n\r\n" +
				`{"jsonrpc":"2.0","method":"notification","params":{"data":1}}`,
			expectError:  false,
			expectID:     nil,
			expectMethod: "notification",
		},
		{
			name: "valid request with string id",
			input: "Content-Length: 67\r\n\r\n" +
				`{"jsonrpc":"2.0","id":"test-id","method":"method","params":{"x":1}}`,
			expectError:  false,
			expectID:     func() *json.RawMessage { raw := json.RawMessage(`"test-id"`); return &raw }(),
			expectMethod: "method",
		},
		{
			name: "invalid json-rpc version",
			input: "Content-Length: 50\r\n\r\n" +
				`{"jsonrpc":"1.0","id":1,"method":"test","params":{}}`,
			expectError: true,
		},
		{
			name: "missing content-length header",
			input: "\r\n" +
				`{"jsonrpc":"2.0","id":1,"method":"test"}`,
			expectError: true,
		},
		{
			name: "invalid content-length",
			input: "Content-Length: invalid\r\n\r\n" +
				`{"jsonrpc":"2.0","id":1,"method":"test"}`,
			expectError: true,
		},
		{
			name: "malformed json",
			input: "Content-Length: 25\r\n\r\n" +
				`{"jsonrpc":"2.0","id":1,}`,
			expectError: true,
		},
		{
			name: "content shorter than declared",
			input: "Content-Length: 100\r\n\r\n" +
				`{"jsonrpc":"2.0","id":1,"method":"test"}`,
			expectError:  false, // Should still parse valid JSON within limit
			expectID:     func() *json.RawMessage { raw := json.RawMessage(`1`); return &raw }(),
			expectMethod: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			request, err := Read(reader)

			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
				return
			}

			require.NoError(t, err, "unexpected error")
			require.NotNil(t, request, "expected request but got nil")

			// Check ID
			if tt.expectID == nil {
				assert.Nil(t, request.ID, "expected no ID but got %v", request.ID)
			} else {
				require.NotNil(t, request.ID, "expected ID %s but got nil", string(*tt.expectID))
				assert.Equal(t, string(*tt.expectID), string(*request.ID), "ID mismatch")
			}

			// Check method
			assert.Equal(t, tt.expectMethod, request.Method, "method mismatch")

			// Check JSON-RPC version
			assert.True(t, request.IsJSONRPC(), "request should have valid JSON-RPC version")
		})
	}
}

func TestProtocolReadEdgeCases(t *testing.T) {
	t.Run("empty reader", func(t *testing.T) {
		reader := bufio.NewReader(strings.NewReader(""))
		_, err := Read(reader)
		assert.Error(t, err, "expected error when reading from empty reader")
	})

	t.Run("only headers no body", func(t *testing.T) {
		input := "Content-Length: 0\r\n\r\n"
		reader := bufio.NewReader(strings.NewReader(input))
		_, err := Read(reader)
		assert.Error(t, err, "expected error when no JSON body provided")
	})

	t.Run("large content length", func(t *testing.T) {
		largeContent := strings.Repeat("x", 1000)
		requestJSON := `{"jsonrpc":"2.0","id":1,"method":"test","params":"` + largeContent + `"}`
		input := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(requestJSON), requestJSON)

		reader := bufio.NewReader(strings.NewReader(input))
		request, err := Read(reader)
		assert.NoError(t, err, "should handle large content")
		assert.NotNil(t, request, "expected request but got nil")
	})
}
