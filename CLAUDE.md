# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

### Core Development
- `make test` - Run all tests
- `make format` - Format code with gofmt and apply license headers
- `make mod` - Tidy Go modules
- `make hygiene` - Run format and mod tidy
- `make all` - Full build pipeline: format, mod, test, and check for dirty working tree
- `make dirty` - Check if working tree is clean (exit code 1 if dirty)
- `make install` - Build and install binary to GOPATH/bin with version info
- `make licenser` - Apply license headers to all files

### Testing Individual Components
- `go test ./pkg/parser/...` - Test parser package
- `go test ./lsp/pkg/...` - Test LSP server packages
- `go test -run TestSpecificFunction ./path/to/package` - Run specific test

### Building
- `go build -o notedown-lsp ./lsp/` - Build LSP server binary
- `go run ./lsp/ serve` - Run LSP server directly
- `make install` - Build with version info and install to GOPATH/bin

### Code Quality
- Uses `licenser` tool for license header management (run `make licenser`)
- All code must be gofmt formatted
- Working tree must be clean after running hygiene tasks

## Architecture Overview

This is a Go-based Language Server Protocol (LSP) implementation for Notedown Flavored Markdown (NFM), consisting of two main components:

### 1. Parser Package (`pkg/parser/`)
- **Core Parser**: Built on `goldmark` with custom extensions for NFM
- **AST Conversion**: Converts goldmark AST to custom tree structure with position tracking
- **Extensions**: Supports wikilinks (`[[target]]` or `[[target|display]]`) via custom goldmark extension
- **Tree Structure**: Implements visitor pattern for AST traversal with precise position information

Key files:
- `parser.go` - Main parser implementation using goldmark
- `tree.go` - Custom AST node types and tree structure
- `visitor.go` - Visitor pattern for tree traversal
- `extensions/wikilink.go` - Wikilink syntax support

### 2. LSP Server (`lsp/`)
- **Server Implementation**: Custom LSP server implementation with JSON-RPC protocol handling
- **JSON-RPC**: Custom JSON-RPC implementation with batch support and error handling in `pkg/jsonrpc/`
- **Command Structure**: Cobra-based CLI with `serve` command for LSP mode
- **Protocol Support**: Full LSP 3.17 specification with comprehensive client/server capabilities

Key components:
- `main.go` + `cmd/` - CLI entry point and commands (serve, version, root)
- `pkg/lsp/` - Core LSP server implementation with mux-based request routing
- `pkg/jsonrpc/` - JSON-RPC protocol handling with batch support, request/response types
- `pkg/constants/` - Server constants and configuration

#### LSP Implementation Details
- **Request Multiplexer**: `mux.go` handles JSON-RPC request routing with structured logging
- **Capability Negotiation**: Comprehensive LSP 3.17 client/server capabilities in separate files:
  - `capabilities_client.go` - Complete client capability structures with LSP spec documentation
  - `capabilities_server.go` - Complete server capability structures with LSP spec documentation
- **Protocol Methods**: All LSP 3.17 methods defined in `methods.go`
- **Initialization**: Proper LSP lifecycle with initialize/initialized sequence

### 3. Shared Packages (`pkg/`)
- **Logging**: Structured logging with multiple levels and formats (`pkg/log/`)
- **Versioning**: Build-time version information (`pkg/version/`)

### Dependencies
- `goldmark` - Markdown parser foundation
- `tliron/glsp` - LSP protocol utilities (used minimally)
- `spf13/cobra` - CLI framework
- `tliron/commonlog` - Structured logging
- Custom wikilink extension for NFM-specific syntax

### Testing Strategy
- Unit tests for parser components (`parser_test.go`)
- JSON-RPC protocol tests with batch handling (`read_test.go`, `write_test.go`)
- Logger tests with different output formats and levels
- Test files use standard Go testing conventions
- Focus on AST conversion accuracy and position tracking

### Language Features
- **Notedown Flavored Markdown**: Opinionated Markdown subset focused on readability
- **Wikilinks**: Internal linking with `[[target]]` or `[[target|display]]` syntax
- **Standard Markdown**: GitHub Flavored Markdown support plus footnotes
- **Semantic Focus**: Emphasizes semantic meaning over HTML rendering

### Development Notes
- All log messages should start with lowercase characters
- Request IDs are properly formatted in logs using `formatRequestID()` helper
- LSP capabilities are split into client (`capabilities_client.go`) and server (`capabilities_server.go`) files
- Uses `any` instead of `interface{}` throughout the codebase
- Comprehensive LSP spec documentation comments on all capability structures

The codebase follows standard Go project structure with clear separation between parsing logic and LSP server functionality. The LSP implementation uses a custom JSON-RPC layer rather than external LSP libraries for maximum control and customization.