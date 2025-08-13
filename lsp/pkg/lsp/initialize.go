package lsp

import "encoding/json"

type InitializeParams struct {
	ClientInfo   *ClientInfo        `json:"clientInfo"`
	Capabilities ClientCapabilities `json:"capabilities"`
}

type InitializeResult struct {
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

func (m *Mux) Run() error {
	m.RegisterMethod(MethodInitialize, func(params json.RawMessage) (any, error) {
		var initParams InitializeParams
		if err := json.Unmarshal(params, &initParams); err != nil {
			m.logger.Error("failed to unmarshal initialize params", "error", err)
			return nil, err
		}

		clientName := "unknown"
		if initParams.ClientInfo != nil {
			clientName = initParams.ClientInfo.Name
		}

		m.logger.Info("lSP client initialized", "client", clientName, "server_version", m.version)

		syncKind := TextDocumentSyncKindFull
		result := InitializeResult{
			ServerInfo: &ServerInfo{Name: "Notedown LSP Server", Version: m.version},
			Capabilities: ServerCapabilities{
				TextDocumentSync: &TextDocumentSyncOptions{
					OpenClose: &[]bool{true}[0],
					Change:    &syncKind,
				},
			},
		}
		return result, nil
	})

	// TODO: Handle initialized message?

	m.logger.Info("starting LSP message processing loop")
	for {
		if err := m.process(); err != nil {
			m.logger.Error("lSP processing error", "error", err)
			return err
		}
	}
}
