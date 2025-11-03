# Notedown Development Guide

Development guide for the Notedown core library and parser.

## Prerequisites

**Option 1: Nix (Recommended)**
- Install [Nix](https://github.com/DeterminateSystems/nix-installer)
- All dependencies provided automatically via `flake.nix`

**Option 2: Manual Setup**
- Go 1.24.4+
- golangci-lint
- licenser
- buf (for protobuf generation)

## Quick Start

```bash
# Clone and test
git clone https://github.com/notedownorg/notedown.git
cd notedown
make test
```

## Development Workflow

### Testing

```bash
make test       # Run all tests
make check      # Full check: generate, format, mod, lint, test
```

### Code Quality

```bash
make hygiene    # Format and tidy
make lint       # Run linter
make all        # Hygiene + test + dirty check
```

### Protobuf Generation

```bash
make generate   # Generate Go code from .proto files
```

## Make Targets

| Target | Description |
|--------|-------------|
| `make test` | Run all tests |
| `make check` | Generate + format + mod + lint + test |
| `make hygiene` | Format and tidy modules |
| `make generate` | Generate protobuf code |
| `make format` | Format code and apply licenses |
| `make mod` | Tidy Go modules |
| `make lint` | Run golangci-lint |
| `make dirty` | Check for uncommitted changes |
| `make all` | Hygiene + test + dirty |

## Architecture

```
├── apis/              # Protobuf API definitions
│   ├── proto/        # .proto files
│   └── go/           # Generated Go code
├── pkg/
│   ├── config/       # Configuration loading
│   ├── log/          # Logging utilities
│   ├── parser/       # Markdown parser with Notedown extensions
│   ├── server/       # Document server
│   └── testdata/     # Test fixtures
└── language/         # Language documentation
```

### Components

- **Parser** (`pkg/parser/`) - Markdown parser with Notedown extensions (wikilinks, tasklists)
- **Config** (`pkg/config/`) - Configuration file discovery and loading
- **Log** (`pkg/log/`) - Structured logging
- **Server** (`pkg/server/`) - Document workspace and filtering

## Related Projects

- [language-server](https://github.com/notedownorg/language-server) - LSP server
- [notedown.nvim](https://github.com/notedownorg/notedown.nvim) - Neovim plugin
