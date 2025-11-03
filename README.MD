# Notedown

Core library and parser for Notedown - a modern, extensible markdown format.

## Features

- **Markdown Parser** - Built on goldmark with Notedown extensions
- **Wikilinks** - Internal linking with `[[target]]` or `[[target|display]]` syntax
- **Task Lists** - Checkbox-based task management
- **Configuration** - Workspace configuration discovery and loading
- **Document Server** - Workspace management and document filtering

## Installation

```bash
go get github.com/notedownorg/notedown
```

## Usage

```go
import "github.com/notedownorg/notedown/pkg/parser"

// Parse markdown
p := parser.New()
doc, err := p.Parse([]byte("# Hello\n\n[[wikilink]] and tasks"))
```

## Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for development instructions.

## Related Projects

- [language-server](https://github.com/notedownorg/language-server) - LSP server implementation
- [notedown.nvim](https://github.com/notedownorg/notedown.nvim) - Neovim plugin

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.
