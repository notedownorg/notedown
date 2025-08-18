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
	"io"
	"net/textproto"
	"strconv"
)

const (
	ContentLengthHeader = "Content-Length"

	// Line endings
	CRLF = "\r\n"
)

func Read(reader *bufio.Reader) (*Request, error) {
	header, err := textproto.NewReader(reader).ReadMIMEHeader()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}
	contentLength, err := strconv.ParseInt(header.Get(ContentLengthHeader), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Length header: %w", err)
	}

	var request Request
	if err := json.NewDecoder(io.LimitReader(reader, contentLength)).Decode(&request); err != nil {
		return nil, fmt.Errorf("error decoding request body: %w", err)
	}
	if !request.IsJSONRPC() {
		return nil, fmt.Errorf("invalid JSON-RPC version: %s", request.ProtocolVersion)
	}
	return &request, nil
}
