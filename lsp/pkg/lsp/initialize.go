package lsp

import "encoding/json"

type InitializeParams struct {
	ClientInfo   *ClientInfo        `json:"clientInfo"`
	Capabilities ClientCapabilities `json:"capabilities"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ClientCapabilities struct {
}

type InitializeResult struct {
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	TextDocumentSync TextDocumentSyncKind `json:"textDocumentSync"`
}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone TextDocumentSyncKind = iota
	TextDocumentSyncKindFull
	TextDocumentSyncKindIncremental
)

func (m *Mux) Run() error {
	m.RegisterMethod(MethodInitialize, func(params json.RawMessage) (any, error) {
		var initParams InitializeParams
		if err := json.Unmarshal(params, &initParams); err != nil {
			return nil, err
		}

		result := InitializeResult{
			ServerInfo: &ServerInfo{Name: "Notedown LSP Server", Version: m.version},
			Capabilities: ServerCapabilities{
				TextDocumentSync: TextDocumentSyncKindFull,
			},
		}
		return result, nil
	})

	// TODO: Handle initialized message?

	for {
		if err := m.process(); err != nil {
			return err
		}
	}
}
