# System Architecture

High-level architecture overview of the Notedown system.

## Components

```
┌─────────────────┐    ┌─────────────────┐
│   LSP Client    │────│   LSP Server    │
│   (Editor)      │    │                 │
└─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │     Parser      │
                       │   (goldmark)    │
                       └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Extensions    │
                       │   (wikilinks)   │
                       └─────────────────┘
```

## Parser Layer

Built on `goldmark` with custom extensions for NFM features:

- **Core Parser**: Standard CommonMark + GFM
- **Wikilink Extension**: `[[target]]` and `[[target|display]]` support
- **AST Conversion**: goldmark AST → custom tree with position tracking

## LSP Layer

Provides language server capabilities:

- **Protocol Handling**: JSON-RPC over stdio
- **Document Management**: Text synchronization and change tracking  
- **Language Features**: Hover, definition, validation

## Data Flow

1. Editor sends document changes via LSP
2. LSP server parses content using [[Parser API|parser]]
3. Wikilinks are resolved and validated
4. Language features provided back to editor

## Key Design Principles

- **Separation of Concerns**: Parser independent of LSP
- **Extensibility**: Plugin-based extension system
- **Performance**: Incremental parsing and caching
- **Standards Compliance**: LSP 3.17 + CommonMark

## Related Files

- [[api-reference]] - Detailed API documentation
- [[notes/technical-decisions]] - Architecture decisions
- [[projects/project-alpha]] - Implementation example