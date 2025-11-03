# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

### Core Development
- `make test` - Run all tests (`go test ./pkg/...`)
- `make format` - Format code with gofmt and apply license headers
- `make mod` - Tidy Go modules
- `make hygiene` - Run format and mod tidy
- `make all` - Full build pipeline: hygiene, test, and check for dirty working tree
- `make dirty` - Check if working tree is clean (exit code 1 if dirty)
- `make check` - Full check: generate, format, mod, lint, test
- `make generate` - Generate Go code from protobuf definitions
- `make lint` - Run golangci-lint
- `make licenser` - Apply license headers to all files

### Testing Individual Components
- `go test ./pkg/parser/...` - Test parser package
- `go test ./pkg/config/...` - Test config package
- `go test ./pkg/server/...` - Test server package
- `go test -run TestSpecificFunction ./path/to/package` - Run specific test

### Code Quality
- Uses `licenser` tool for license header management (run `make licenser`)
- All Go code must be gofmt formatted
- Working tree must be clean after running hygiene tasks
- Requires Go 1.24.4 or later
- Optional Nix development environment available (flake.nix)

## Architecture Overview

This is the core Notedown library providing the parser, configuration, and shared utilities.

### 1. Parser Package (`pkg/parser/`)
- **Core Parser**: Built on `goldmark` with custom extensions for Notedown Flavored Markdown
- **AST Conversion**: Converts goldmark AST to custom tree structure with position tracking
- **Extensions**: Supports wikilinks (`[[target]]` or `[[target|display]]`) and tasklists via custom goldmark extensions
- **Tree Structure**: Implements visitor pattern for AST traversal with precise position information

Key files:
- `parser.go` - Main parser implementation using goldmark
- `tree.go` - Custom AST node types and tree structure
- `visitor.go` - Visitor pattern for tree traversal
- `extensions/wikilink.go` - Wikilink syntax support
- `extensions/tasklist.go` - Task list syntax support

### 2. Configuration Package (`pkg/config/`)
- **Discovery**: Automatic configuration file discovery in workspace
- **Loading**: YAML configuration file parsing
- **Types**: Configuration type definitions and validation

### 3. Logging Package (`pkg/log/`)
- **Structured Logging**: Multiple log levels (debug, info, warn, error)
- **Formats**: Text and JSON output formats
- **LSP Support**: Special LSP-compatible logging mode

### 4. Server Package (`pkg/server/`)
- **Workspace**: Document workspace management
- **Filtering**: Document filtering and discovery
- **Loading**: Document content loading

### Dependencies
- `goldmark` - Markdown parser foundation
- `spf13/cobra` - CLI framework (for version info)
- `stretchr/testify` - Testing framework
- `go.abhg.dev/goldmark/frontmatter` - Frontmatter support

### Testing Strategy
- **Unit Tests**: Parser components (`parser_test.go`), extension tests
- **Integration Tests**: Configuration loading, workspace management
- **Conventions**: Standard Go testing with table-driven tests where appropriate

### Language Features
- **Notedown Flavored Markdown**: Opinionated Markdown subset focused on readability and semantic meaning
- **Wikilinks**: Internal linking with `[[target]]` or `[[target|display]]` syntax
- **Task Lists**: Checkbox-based task management
- **Standard Markdown**: GitHub Flavored Markdown support plus footnotes
- **Semantic Focus**: Emphasizes semantic meaning over HTML rendering
- **Language Specification**: Full language documentation available in `language/` directory

### Related Projects
- **Language Server**: LSP implementation in separate `language-server` repository
- **Neovim Plugin**: Editor integration in separate `notedown.nvim` repository

### Development Notes
- All log messages should start with lowercase characters
- Uses `any` instead of `interface{}` throughout the codebase
- Configuration files use YAML format (`.notedown/settings.yaml`)

The codebase follows standard Go project structure with clear separation of concerns across packages.
