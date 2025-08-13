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
