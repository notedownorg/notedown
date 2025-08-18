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
	"encoding/json"
)

const JSONRPCVersion = "2.0"

type Message interface {
	IsJSONRPC() bool
}

// Request represents a JSON-RPC 2.0 request
type Request struct {
	ProtocolVersion string           `json:"jsonrpc"`
	ID              *json.RawMessage `json:"id"`
	Method          string           `json:"method"`
	Params          json.RawMessage  `json:"params"`
}

func (r Request) IsJSONRPC() bool {
	return r.ProtocolVersion == JSONRPCVersion
}

func (r Request) IsNotification() bool {
	return r.ID == nil
}

// NewNotification creates a new JSON-RPC notification (no ID)
func NewNotification(method string, params any) *Request {
	var paramsBytes json.RawMessage
	if params != nil {
		if b, err := json.Marshal(params); err == nil {
			paramsBytes = b
		}
	}

	return &Request{
		ProtocolVersion: JSONRPCVersion,
		ID:              nil, // Notifications have no ID
		Method:          method,
		Params:          paramsBytes,
	}
}
