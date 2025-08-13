package lsp

type InitializeParams struct {
	ClientInfo   *ClientInfo        `json:"clientInfo"`
	Capabilities ClientCapabilities `json:"capabilities"`
}

type InitializeResult struct {
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
	Capabilities ServerCapabilities `json:"capabilities"`
}
