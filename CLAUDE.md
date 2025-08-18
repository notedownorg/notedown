# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

### Core Development
- `make test` - Run all tests (LSP, Neovim, and integration)
- `make test-lsp` - Run Go tests only (`go test ./...`)
- `make test-nvim` - Run Neovim plugin tests only
- `make format` - Format code with gofmt
- `make mod` - Tidy Go modules
- `make hygiene` - Run format and mod tidy
- `make all` - Full build pipeline: format, mod, test, and check for dirty working tree
- `make dirty` - Check if working tree is clean (exit code 1 if dirty)
- `make install` - Build and install `notedown-language-server` binary to GOPATH/bin with version info, also installs Neovim plugin to `~/.config/notedown/nvim/`
- `make clean` - Remove installed binary and Neovim plugin files
- `make licenser` - Apply license headers to all files
- `make dev` - Install and open test workspace in Neovim for development

### Testing Individual Components
- `go test ./pkg/parser/...` - Test parser package
- `go test ./lsp/pkg/...` - Test LSP server packages
- `go test ./lsp/pkg/notedownls/...` - Test Notedown-specific LSP implementation
- `go test -run TestSpecificFunction ./path/to/package` - Run specific test
- `cd notedown-nvim && nvim --headless --noplugin -u tests/helpers/minimal_init.lua -c "lua MiniTest.run()" -c "qall!"` - Test Neovim plugin

### Building
- `go build -o bin/notedown-language-server ./lsp/` - Build LSP server binary
- `go run ./lsp/ serve` - Run LSP server directly
- `make install` - Build with version info and install to GOPATH/bin

### Code Quality
- Uses `licenser` tool for license header management (run `make licenser`)
- All code must be gofmt formatted (Go) and stylua formatted (Lua)
- Working tree must be clean after running hygiene tasks
- Requires Go 1.24.4 or later
- Optional Nix development environment available (flake.nix)

## Architecture Overview

This is a Go-based Language Server Protocol (LSP) implementation for Notedown Flavored Markdown (NFM), consisting of four main components:

### 1. Parser Package (`pkg/parser/`)
- **Core Parser**: Built on `goldmark` with custom extensions for NFM
- **AST Conversion**: Converts goldmark AST to custom tree structure with position tracking
- **Extensions**: Supports wikilinks (`[[target]]` or `[[target|display]]`) and tasklists via custom goldmark extensions
- **Tree Structure**: Implements visitor pattern for AST traversal with precise position information

Key files:
- `parser.go` - Main parser implementation using goldmark
- `tree.go` - Custom AST node types and tree structure
- `visitor.go` - Visitor pattern for tree traversal
- `extensions/wikilink.go` - Wikilink syntax support
- `extensions/tasklist.go` - Task list syntax support

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
- `pkg/notedownls/` - Notedown-specific LSP server implementation

#### LSP Implementation Details
- **Request Multiplexer**: `mux.go` handles JSON-RPC request routing with structured logging
- **Capability Negotiation**: Comprehensive LSP 3.17 client/server capabilities in separate files:
  - `capabilities_client.go` - Complete client capability structures with LSP spec documentation
  - `capabilities_server.go` - Complete server capability structures with LSP spec documentation
- **Protocol Methods**: All LSP 3.17 methods defined in `methods.go`
- **Initialization**: Proper LSP lifecycle with initialize/initialized sequence

#### Notedown Language Server Implementation (`lsp/pkg/notedownls/`)
- **Server**: Main notedownls Server struct that implements LSP interface for Notedown features
- **Document Management**: Thread-safe document storage with content tracking and lifecycle management
- **Workspace Management**: WorkspaceManager handles workspace roots, file discovery, and Markdown file indexing
- **Wikilink Features**: 
  - Completion provider with intelligent target suggestions (existing files, referenced targets, directory paths)
  - Definition provider for wikilink navigation
  - Context-aware parsing of wikilink syntax with pipe separator support (`[[target|display]]`)
- **Indexing System**: WikilinkIndex tracks all wikilink targets across workspace with reference counting

Key files in notedownls:
- `server.go` - Main server implementation with initialization and handler registration
- `document.go` - Document lifecycle management and content tracking
- `workspace.go` - Workspace discovery and file management
- `handlers_textDocument.go` - LSP text document method handlers (completion, definition)
- `handlers_workspace.go` - LSP workspace method handlers
- `indexes/wikilink.go` - Advanced wikilink target indexing and resolution system

### 3. Neovim Plugin (`notedown-nvim/`)
- **Plugin System**: Lua-based Neovim plugin with intelligent workspace detection
- **LSP Integration**: Automatic LSP server startup and configuration
- **Parser Modes**: Configurable parser selection (auto/notedown/markdown)
- **Workspace Detection**: Path-based workspace matching with tilde expansion

Key files:
- `lua/notedown/init.lua` - Main plugin initialization and configuration
- `lua/notedown/config.lua` - Configuration management and workspace detection
- `plugin/notedown.lua` - Plugin bootstrapping
- `tests/` - Mini.nvim-based test suite

### 4. Shared Packages (`pkg/`)
- **Logging**: Structured logging with multiple levels and formats (`pkg/log/`)
- **Versioning**: Build-time version information (`pkg/version/`)

### Dependencies
- `goldmark` - Markdown parser foundation
- `spf13/cobra` - CLI framework
- `stretchr/testify` - Testing framework
- Custom wikilink and tasklist extensions for NFM-specific syntax

### Testing Strategy
- Unit tests for parser components (`parser_test.go`)
- JSON-RPC protocol tests with batch handling (`read_test.go`, `write_test.go`)
- Logger tests with different output formats and levels
- Notedownls tests cover completion, workspace management, and wikilink indexing
- Neovim plugin tests with Mini.nvim framework for realistic scenarios
- Demo workspace available in `demo_workspace/` with sample Markdown files and directory structure, including examples of ambiguous wikilinks
- Test files use standard Go testing conventions
- Focus on AST conversion accuracy and position tracking

### Wikilink Index System
The `indexes/wikilink.go` implements a sophisticated wikilink tracking system:
- **Target Tracking**: Maintains registry of all wikilink targets with existence status and reference counts
- **Reference Management**: Tracks which documents reference each target, enabling cleanup of unused targets
- **Existence Detection**: Matches wikilink targets to actual workspace files (both filename and full path matching)
- **Completion Intelligence**: Provides three tiers of completion suggestions:
  1. Existing files (highest priority)
  2. Non-existent targets referenced elsewhere (medium priority)  
  3. Directory path completions (lowest priority)
- **Thread Safety**: All operations are mutex-protected for concurrent access
- **Regex-based Extraction**: Uses regex parsing for reliable wikilink detection in document content

### Language Features
- **Notedown Flavored Markdown**: Opinionated Markdown subset focused on readability and semantic meaning
- **Wikilinks**: Internal linking with `[[target]]` or `[[target|display]]` syntax
- **Task Lists**: Checkbox-based task management
- **Standard Markdown**: GitHub Flavored Markdown support plus footnotes  
- **Semantic Focus**: Emphasizes semantic meaning over HTML rendering
- **Language Specification**: Full language documentation available in `language/` directory

### Editor Integration
- **Neovim Plugin**: Complete Neovim integration provided in `notedown-nvim/`
  - Lua-based configuration and initialization (`lua/notedown/`)
  - Plugin bootstrapping (`plugin/notedown.lua`)
  - Automatically installed to `~/.config/notedown/nvim/` with `make install`
- **LSP Integration**: Compatible with any LSP-capable editor through standard LSP protocol

### Development Notes
- All log messages should start with lowercase characters
- Request IDs are properly formatted in logs using `formatRequestID()` helper
- LSP capabilities are split into client (`capabilities_client.go`) and server (`capabilities_server.go`) files
- Uses `any` instead of `interface{}` throughout the codebase
- Comprehensive LSP spec documentation comments on all capability structures

The codebase follows standard Go project structure with clear separation between parsing logic and LSP server functionality. The LSP implementation uses a custom JSON-RPC layer rather than external LSP libraries for maximum control and customization.
