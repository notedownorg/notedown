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
	"fmt"
)

type Response struct {
	ProtocolVersion string           `json:"jsonrpc"`
	ID              *json.RawMessage `json:"id"`
	Result          any              `json:"result,omitempty"`
	Error           *ResponseError   `json:"error,omitempty"`
}

func (r Response) IsJSONRPC() bool {
	return r.ProtocolVersion == "2.0"
}

// NewResponse creates a successful response
func NewResponse(id *json.RawMessage, result any) Response {
	return Response{
		ProtocolVersion: JSONRPCVersion,
		ID:              id,
		Result:          result,
	}
}

// ResponseError represents a JSON-RPC 2.0 error object
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Error implements the error interface
func (e *ResponseError) Error() string {
	return e.Message
}

// Standard JSON-RPC 2.0 error codes
const (
	ParseErrorCode     = -32700
	InvalidRequestCode = -32600
	MethodNotFoundCode = -32601
	InvalidParamsCode  = -32602
	InternalErrorCode  = -32603

	// Reserved error code range
	ReservedErrorMin = -32768
	ReservedErrorMax = -32000
)

// NewParseError creates a parse error response
func NewParseError(id *json.RawMessage, data any) Response {
	return Response{
		ProtocolVersion: JSONRPCVersion,
		ID:              id,
		Error: &ResponseError{
			Code:    ParseErrorCode,
			Message: "Parse error",
			Data:    data,
		},
	}
}

// NewInvalidRequestError creates an invalid request error response
func NewInvalidRequestError(id *json.RawMessage, data any) Response {
	return Response{
		ProtocolVersion: JSONRPCVersion,
		ID:              id,
		Error: &ResponseError{
			Code:    InvalidRequestCode,
			Message: "Invalid Request",
			Data:    data,
		},
	}
}

// NewMethodNotFoundError creates a method not found error response
func NewMethodNotFoundError(id *json.RawMessage, data any) Response {
	return Response{
		ProtocolVersion: JSONRPCVersion,
		ID:              id,
		Error: &ResponseError{
			Code:    MethodNotFoundCode,
			Message: "Method not found",
			Data:    data,
		},
	}
}

// NewInvalidParamsError creates an invalid params error response
func NewInvalidParamsError(id *json.RawMessage, data any) Response {
	return Response{
		ProtocolVersion: JSONRPCVersion,
		ID:              id,
		Error: &ResponseError{
			Code:    InvalidParamsCode,
			Message: "Invalid params",
			Data:    data,
		},
	}
}

// NewInternalError creates an internal error response
func NewInternalError(id *json.RawMessage, data any) Response {
	return Response{
		ProtocolVersion: JSONRPCVersion,
		ID:              id,
		Error: &ResponseError{
			Code:    InternalErrorCode,
			Message: "Internal error",
			Data:    data,
		},
	}
}

// NewCustomError creates a custom error response with code validation
func NewCustomError(id *json.RawMessage, code int, message string, data any) (Response, error) {
	if err := validateErrorCode(code); err != nil {
		return Response{}, err
	}
	return Response{
		ProtocolVersion: JSONRPCVersion,
		ID:              id,
		Error: &ResponseError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}, nil
}

// validateErrorCode checks if an error code is in the reserved range
func validateErrorCode(code int) error {
	if code >= ReservedErrorMin && code <= ReservedErrorMax {
		// Allow standard JSON-RPC error codes
		switch code {
		case ParseErrorCode, InvalidRequestCode, MethodNotFoundCode, InvalidParamsCode, InternalErrorCode:
			return nil
		default:
			return fmt.Errorf("error code %d is in the reserved range %d to %d", code, ReservedErrorMin, ReservedErrorMax)
		}
	}
	return nil
}

// IsParseError checks if an error is a parse error
func (e *ResponseError) IsParseError() bool {
	return e.Code == ParseErrorCode
}

// IsInvalidRequestError checks if an error is an invalid request error
func (e *ResponseError) IsInvalidRequestError() bool {
	return e.Code == InvalidRequestCode
}

// IsMethodNotFoundError checks if an error is a method not found error
func (e *ResponseError) IsMethodNotFoundError() bool {
	return e.Code == MethodNotFoundCode
}

// IsInvalidParamsError checks if an error is an invalid params error
func (e *ResponseError) IsInvalidParamsError() bool {
	return e.Code == InvalidParamsCode
}

// IsInternalError checks if an error is an internal error
func (e *ResponseError) IsInternalError() bool {
	return e.Code == InternalErrorCode
}

// IsStandardError checks if an error is one of the standard JSON-RPC errors
func (e *ResponseError) IsStandardError() bool {
	return e.IsParseError() || e.IsInvalidRequestError() || e.IsMethodNotFoundError() ||
		e.IsInvalidParamsError() || e.IsInternalError()
}

// IsReservedError checks if an error code is in the reserved range
func (e *ResponseError) IsReservedError() bool {
	return e.Code >= ReservedErrorMin && e.Code <= ReservedErrorMax
}
