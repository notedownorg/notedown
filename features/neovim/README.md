# Notedown Neovim Features

This directory contains **living documentation** for Notedown's Neovim plugin features. Each feature is documented through executable VHS (Video Home System) tests that serve as both verification and demonstration.

## Living Documentation Philosophy

Rather than static documentation that becomes outdated, we use **executable documentation** that:
- **Proves** features work by running them
- **Demonstrates** features visually through generated GIFs
- **Validates** behavior through automated regression testing
- **Documents** usage patterns through real terminal interactions

## Areas

### [Initialization](./initialization/)
Plugin loading, setup, and workspace detection features.

- **[Workspace Status Command](./initialization/workspace-status-command/)** - `:NotedownWorkspaceStatus` command

## Directory Structure

```
features/neovim/
├── README.md                           # This file - feature overview
├── framework_test.go                   # Test runner for all features
├── pkg/notedown/runner.go              # Clean testing framework
├── initialization/                     # AREA: Plugin setup and startup
│   ├── README.md                      # Area documentation with features
│   └── workspace-status-command/       # FEATURE: :NotedownWorkspaceStatus
│       ├── README.md                  # Feature documentation
│       ├── demo.gif                   # Visual demonstration
│       ├── demo.tape.tmpl             # VHS template
│       ├── workspace/                 # Test workspace
│       └── expected.ascii             # Expected output
└── [future areas]/
    └── [features]/
```

## Running Features

### All Features
```bash
# Run all feature tests
make test-features

# Run specific area/feature
go test -run TestFeatures/initialization/workspace-status-command
```

### Individual Features
Each feature directory contains a demo.gif that shows the feature in action. To regenerate:
```bash
cd features/neovim/initialization/workspace-status-command
vhs demo.tape.tmpl
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