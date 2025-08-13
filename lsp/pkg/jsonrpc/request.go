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
