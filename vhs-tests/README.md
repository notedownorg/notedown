# VHS Tests for Notedown Neovim Plugin

This directory contains VHS (Video Home System) tests for the Notedown Neovim plugin, providing end-to-end terminal-based testing that complements the existing mini.test unit tests.

## Directory Structure

- `testdata/templates/` - VHS template files (`.tape.tmpl`) for dynamic test generation
- `golden/` - Expected ASCII output files for regression testing
- `gifs/` - Generated GIF files for visual inspection (not committed)
- `workspaces/` - Test workspace directories for isolated testing scenarios
- `pkg/notedown/` - Clean VHS testing framework with `NotedownVHSRunner`
- `framework_test.go` - Test definitions and execution logic

## Test Framework Architecture

### Clean Runner Design
- **NotedownVHSRunner** - Handles all testing boilerplate (LSP building, plugin installation, workspace setup, cleanup)
- **Data-Driven Tests** - Tests defined as simple `VHSTest` structs in `framework_test.go`
- **Template Rendering** - Dynamic VHS tape generation from Go templates with test-specific data
- **Parallel Execution** - Safe concurrent test execution with shared LSP binary building via `sync.Once`

### Test Categories
- **Plugin Initialization** - Test automatic plugin loading and LSP startup
- **Wikilink Navigation** - Test go-to-definition functionality  
- **Wikilink Completion** - Test autocompletion in real editing scenarios
- **Wikilink Diagnostics** - Test ambiguous wikilink warnings and display

## Running Tests

### Make Targets
```bash
# Run all VHS tests
make test-vhs

# Run all tests including VHS
make test
```

### Go Test Usage
```bash
# Run all VHS tests with parallel execution
go test -parallel 4 -v -timeout 20m

# Run specific test
go test -v -run TestVHSFramework/plugin-initialization

# Run with increased parallelism
go test -parallel 8 -v
```

## Test Approach

### Golden File Testing
- Uses VHS `.ascii` output format for deterministic comparison
- Compares terminal output against golden files with testify assertions
- **Golden files auto-created** if missing on first test run
- **Golden files are committed to git** for regression detection
- Enables automated CI testing without visual dependencies

### Visual Output
- Generates GIFs automatically in `gifs/` directory for manual inspection
- Useful for debugging and documentation (not committed to git)
- Both ASCII and GIF outputs generated simultaneously
- **Note**: GIF generation requires a graphics environment - files may be empty (0 bytes) in headless CI environments

### Parallel Execution
- Tests run concurrently for faster feedback (4 parallel by default)
- Shared LSP binary building with `sync.Once` avoids redundant compilation
- Per-test cleanup and isolated temporary directories
- Increased timeouts (300s) for parallel execution stability

## Template System

### Dynamic VHS Generation
Templates in `testdata/templates/` use Go template syntax:
```vhs
Output "{{.OutputFile}}"
Output "{{.TmpDir}}/{{.TestName}}.gif"

Set WorkingDirectory "{{.WorkspaceDir}}"
Set FontSize 16
Set Width 1200
Set Height 800

Type "nvim --noplugin -u {{.ConfigFile}} index.md"
Enter
Sleep 3s

Type ":LspInfo"
Enter
Sleep 2s
```

### Available Template Data
- `OutputFile` - Path for ASCII output
- `WorkspaceDir` - Test workspace directory
- `ConfigFile` - Neovim configuration file path
- `TmpDir` - Temporary directory for test files
- `LSPBinary` - Path to built LSP server

## Available Tests

Current test suite includes:
- `plugin-initialization` - Plugin loading and LSP startup (workspace: plugin-init-test)
- `wikilink-navigation` - Go-to-definition functionality (workspace: plugin-init-test)  
- `wikilink-completion` - Autocompletion for wikilinks (workspace: completion-test)
- `wikilink-diagnostics` - Ambiguous wikilink warnings (workspace: diagnostics-test)

## Writing New Tests

1. **Add test definition** to `vhsTests` slice in `framework_test.go`:
   ```go
   {Name: "new-feature", Workspace: "test-workspace", Timeout: 300 * time.Second}
   ```

2. **Create template** in `testdata/templates/new-feature.tape.tmpl`:
   ```vhs
   Output "{{.OutputFile}}"
   Output "{{.TmpDir}}/new-feature.gif"
   
   Set WorkingDirectory "{{.WorkspaceDir}}"
   # ... VHS commands for your test
   ```

3. **Create workspace** in `workspaces/test-workspace/` with test files

4. **Run test** to generate golden file: `go test -v -run TestVHSFramework/new-feature`

5. **Verify output** by checking generated golden file and GIF

## Dependencies

- VHS (Video Home System) by Charmbracelet (available in nix flake)
- notedown-language-server (built automatically from parent directory)
- Neovim with notedown plugin (installed automatically during test)
- Go templates for dynamic VHS tape generation

## Integration with CI

VHS tests are integrated into GitHub Actions workflow:
- Runs in parallel with unit tests and language server tests
- Uses existing nix development environment (VHS pre-installed)
- Generates both ASCII regression tests and visual GIF artifacts

## Framework Benefits

The refactored framework provides:
- **Clean separation** - Test logic separated from boilerplate
- **Easy maintenance** - New tests require minimal code
- **Fast execution** - Parallel testing with shared binary building  
- **Visual debugging** - Automatic GIF generation for inspection
- **Reliable CI** - Deterministic ASCII output for regression testing