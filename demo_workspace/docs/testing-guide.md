# Testing Guide

Comprehensive testing strategies and guidelines.

## Test Structure

### Unit Tests
- Parser functionality: `pkg/parser/parser_test.go`
- Wikilink extension: `pkg/parser/extensions/wikilink_test.go`
- JSON-RPC protocol: `language-server/pkg/jsonrpc/*_test.go`

### Integration Tests
- End-to-end LSP communication
- Document parsing workflows
- Cross-reference validation

## Testing Patterns

### Table-Driven Tests

```go
func TestWikilinks(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected []Wikilink
    }{
        {
            name:  "simple wikilink",
            input: "[[target]]",
            expected: []Wikilink{{Target: "target"}},
        },
        {
            name:  "wikilink with display",
            input: "[[target|display]]",
            expected: []Wikilink{{Target: "target", Display: "display"}},
        },
    }
    // Test implementation...
}
```

### Mock Objects
- Use interfaces for testability
- Mock file system operations
- Mock LSP client communication

## Running Tests

```bash
# All tests
make test

# Specific package
go test ./pkg/parser/...

# With coverage
go test -cover ./...

# Specific test
go test -run TestWikilinks ./pkg/parser/extensions/
```

## Test Data

- Use `testdata/` directories for fixtures
- Golden files for expected outputs
- Minimal, focused test cases

## Coverage Goals

| Component | Target Coverage |
|-----------|----------------|
| Parser Core | 90%+ |
| Extensions | 85%+ |
| LSP Server | 80%+ |
| Utilities | 95%+ |

## Related Documentation

- [[api-reference]] - API being tested
- [[notes/technical-decisions]] - Testing approach rationale
- [[projects/project-alpha]] - Example test implementations