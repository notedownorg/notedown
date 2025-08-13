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

### Testing Individual Components
- `go test ./pkg/parser/...` - Test parser package
- `go test ./lsp/pkg/...` - Test LSP server packages
- `go test -run TestSpecificFunction ./path/to/package` - Run specific test

### Building
- `go build -o notedown-lsp ./lsp/` - Build LSP server binary
- `go run ./lsp/ serve` - Run LSP server directly

### Code Quality
- Uses `licenser` tool for license header management
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
- **Server Implementation**: Uses `tliron/glsp` library for LSP protocol handling
- **JSON-RPC**: Custom JSON-RPC implementation with batch support and error handling
- **Command Structure**: Cobra-based CLI with `serve` command for LSP mode
- **Protocol Support**: LSP 3.17 with fallback to 3.16 features

Key components:
- `main.go` + `cmd/` - CLI entry point and commands
- `pkg/server/` - LSP server implementation
- `pkg/jsonrpc/` - JSON-RPC protocol handling with batch support
- `pkg/constants/` - Server constants and configuration

### Dependencies
- `goldmark` - Markdown parser foundation
- `tliron/glsp` - LSP protocol implementation
- `spf13/cobra` - CLI framework
- Custom wikilink extension for NFM-specific syntax

### Testing Strategy
- Unit tests for parser components (`parser_test.go`)
- JSON-RPC protocol tests with batch handling
- Test files use standard Go testing conventions
- Focus on AST conversion accuracy and position tracking

### Language Features
- **Notedown Flavored Markdown**: Opinionated Markdown subset focused on readability
- **Wikilinks**: Internal linking with `[[target]]` or `[[target|display]]` syntax
- **Standard Markdown**: GitHub Flavored Markdown support plus footnotes
- **Semantic Focus**: Emphasizes semantic meaning over HTML rendering

The codebase follows standard Go project structure with clear separation between parsing logic and LSP server functionality.