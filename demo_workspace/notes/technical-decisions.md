# Technical Decisions

Documentation of key technical decisions and their rationale.

## Architecture Decisions

### Wikilink Syntax
**Decision**: Use `[[target]]` and `[[target|display]]` syntax for internal links.

**Rationale**: 
- Familiar to users of other wiki systems
- Clean, readable syntax
- Easy to parse and validate

**Related**: See implementation in [[projects/project-beta]]

### Parser Foundation
**Decision**: Build on top of `goldmark` parser.

**Rationale**:
- Mature, well-tested CommonMark implementation
- Extensible architecture
- Good performance characteristics

### LSP Integration
**Decision**: Use `tliron/glsp` for LSP protocol handling.

**Benefits**:
- Handles LSP protocol details
- JSON-RPC support
- Active maintenance

## Code Standards

| Area | Standard | Notes |
|------|----------|-------|
| Formatting | `gofmt` | Enforced by CI |
| Testing | Table-driven | See [[docs/testing-guide]] |
| Documentation | Inline + wiki | Use [[Wikilinks]] liberally |

## Future Considerations

- Performance optimization for large documents
- Multi-language support
- Plugin architecture

Related: [[notes/ideas]] for future enhancements.