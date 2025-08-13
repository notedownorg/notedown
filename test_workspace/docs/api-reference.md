# API Reference

Complete API documentation for the Notedown system.

## Parser API

### `parser.New(options ...Option) *Parser`

Creates a new parser instance with the specified options.

```go
p := parser.New(
    parser.WithWikilinks(true),
    parser.WithExtensions(extensions.Table, extensions.Strikethrough),
)
```

### `Parser.Parse(source []byte) (Node, error)`

Parses markdown source and returns the AST root node.

**Example**:
```go
content := []byte("# Hello [[world]]")
node, err := p.Parse(content)
```

## Wikilink Extension

Supports both simple and display wikilinks:
- `[[target]]` - Simple link
- `[[target|display text]]` - Link with custom display

## LSP Server

### Initialization

The LSP server supports the following capabilities:
- Text document synchronization
- Hover information for [[Wikilinks]]
- Go-to-definition for internal links

### Configuration

Server can be configured via initialization parameters:

```json
{
    "workspaceFolder": "/path/to/workspace",
    "features": {
        "wikilinks": true,
        "validation": true
    }
}
```

## Related Documentation

- [[docs/architecture]] - System design
- [[notes/technical-decisions]] - Implementation rationale
- [[projects/project-alpha]] - Usage examples