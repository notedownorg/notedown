# Feature Development Guide

This guide covers how to develop, test, and maintain the Notedown Neovim plugin feature documentation system.

## Architecture

### Directory Structure

```
features/neovim/
├── README.md                           # User-focused feature overview
├── DEVELOPMENT.md                      # This file - development guide
├── framework_test.go                   # Test runner for all features
├── pkg/notedown/runner.go              # Clean testing framework
├── {area}/                             # AREA: Functional grouping
│   ├── README.md                      # Area documentation
│   └── {feature}/                      # FEATURE: Specific functionality
│       ├── README.md                  # Feature documentation
│       ├── demo.gif                   # Visual demonstration
│       ├── demo.tape.tmpl             # VHS template
│       ├── workspace/                 # Test workspace
│       └── expected.ascii             # Expected output
└── [future areas]/
    └── [features]/
```

### Test Framework

The `NotedownVHSRunner` provides clean separation between test logic and boilerplate:

- **LSP Building**: Shared binary compilation with `sync.Once`
- **Plugin Installation**: Isolated Neovim plugin setup per test
- **Workspace Creation**: Dynamic test workspace generation
- **Template Rendering**: Go template-based VHS tape generation
- **Parallel Execution**: Concurrent test runs with proper isolation

## Running Tests

### All Features
```bash
# Run all feature tests
make test-features

# Run with fresh golden files
make test-features-golden

# Run specific area/feature
go test -v -run TestFeatures/initialization/workspace-status-command
```

### Development Workflow
```bash
# Run a single test during development
cd features/neovim
go test -v -run TestFeatures/initialization/workspace-status-command

# Regenerate golden file (delete and re-run)
rm initialization/workspace-status-command/expected.ascii
go test -v -run TestFeatures/initialization/workspace-status-command
```

## Testing Approach

### Golden File Testing
- Uses VHS `.ascii` output format for deterministic comparison
- Compares terminal output against golden files with testify assertions
- **Golden files auto-created** if missing on first test run
- **Golden files are committed to git** for regression detection
- Enables automated CI testing without visual dependencies

### Visual Output
- Generates GIFs automatically for manual inspection
- Both ASCII and GIF outputs generated simultaneously
- **Note**: GIF generation requires a graphics environment

### Parallel Execution
- Tests run concurrently for faster feedback
- Shared LSP binary building avoids redundant compilation
- Per-test cleanup and isolated temporary directories
- Increased timeouts (300s) for stability

## Template System

### Dynamic VHS Generation
Templates use Go template syntax with available data:

```vhs
Output "{{.OutputFile}}"
Output "{{.TmpDir}}/{{.TestName}}.gif"
Set Width 1200
Set Height 800
Set TypingSpeed 50ms

# Change to test workspace
Type "cd '{{.WorkspaceDir}}'"
Enter

# Start Neovim with plugin
Type "nvim --clean -u '{{.ConfigFile}}' ."
Enter
Sleep 4s

# Test feature
Type ":NotedownWorkspaceStatus"
Enter
Sleep 3s
```

### Available Template Data
- `OutputFile` - Path for ASCII output
- `WorkspaceDir` - Test workspace directory  
- `ConfigFile` - Neovim configuration file path
- `TmpDir` - Temporary directory for test files
- `LSPBinary` - Path to built LSP server
- `TestName` - Safe test name for file operations

## Writing New Features

### 1. Plan the Feature Structure

Determine the area and feature name:
- Area: Functional grouping (e.g., "initialization", "editing", "navigation")  
- Feature: Specific functionality (e.g., "workspace-status-command", "auto-complete")

### 2. Create Directory Structure

```bash
mkdir -p features/neovim/{area}/{feature}/workspace
```

### 3. Add Test Definition

Add to `featureTests` in `framework_test.go`:
```go
{Area: "area-name", Feature: "feature-name", Workspace: "workspace", Timeout: 300 * time.Second},
```

### 4. Create VHS Template

Create `{area}/{feature}/demo.tape.tmpl`:
```vhs
# VHS tape for testing {feature} functionality
Output "{{.OutputFile}}"
Output "{{.TmpDir}}/{{.TestName}}.gif"
Set FontSize 12
Set Width 1200  
Set Height 800
Set TypingSpeed 50ms

# Hide initial setup
Hide

# Setup workspace
Type "cd '{{.WorkspaceDir}}'"
Enter
Sleep 1s

# Start Neovim
Type "nvim --clean -u '{{.ConfigFile}}' index.md"
Enter
Sleep 4s

# Show recording
Show

# Test your feature here
Type ":YourCommand"
Enter
Sleep 3s

# Hide cleanup
Hide
Type ":qa!"
Enter
```

### 5. Create Test Workspace

Create files in `{area}/{feature}/workspace/`:
- Add `.notedown/.gitkeep` for Notedown workspace detection
- Add test files (e.g., `index.md`) with relevant content
- Include any special test cases your feature requires

### 6. Run Initial Test

```bash
go test -v -run TestFeatures/{area}/{feature}
```

This will:
- Build the LSP server
- Install the plugin
- Run your VHS template  
- Generate the golden file
- Create the demo GIF

### 7. Create Documentation

Create `{area}/{feature}/README.md`:
```markdown
# Feature Name

Brief description of what this feature does.

## Demo

![Feature Demo](./demo.gif)

## Usage

How to use this feature:
1. Step one
2. Step two
3. Step three

## Example

Show example usage with code blocks or screenshots.
```

### 8. Update Area Documentation  

Update `{area}/README.md` to include your new feature.

## Framework Internals

### NotedownVHSRunner

The test runner handles:

1. **LSP Binary Building**: Shared compilation using `sync.Once`
2. **Plugin Installation**: Copies Neovim plugin files with correct directory detection
3. **Workspace Creation**: Sets up isolated test workspaces
4. **Template Rendering**: Processes Go templates with test data
5. **VHS Execution**: Runs VHS with proper timeout handling
6. **Asset Management**: Copies GIFs and manages golden files

### Plugin Installation Process

The runner:
1. Finds the project root containing the `neovim/` directory
2. Copies all plugin files to an isolated test directory
3. Creates a custom `init.lua` with proper LSP configuration
4. Verifies plugin files copied correctly

### Template Processing

Templates are processed with:
1. Safe filename generation (replacing `/` with `-`)
2. Proper path resolution for area/feature structure
3. Dynamic data injection for workspace and binary paths
4. Error handling for missing templates

## CI Integration

### GitHub Actions

Tests run in the `features` job in `.github/workflows/test.yml`:
- Uses Nix development environment (VHS pre-installed)
- Runs in parallel with other test suites
- Generates ASCII output for regression testing
- GIFs may be empty in headless CI environments

### Quality Checks

Before committing:
```bash
make format  # Code formatting
make mod     # Go module tidying
```

## Troubleshooting

### Common Issues

1. **Plugin not loading**: Check that `neovim/lua/notedown/config.lua` exists
2. **Template errors**: Verify all template variables are available
3. **Golden file mismatches**: Delete and regenerate golden files
4. **VHS timeouts**: Increase timeout in test definition
5. **Path issues**: Use `safeName` for file operations, not raw test names

### Debug Mode

Add debug output to runner:
```go
t.Logf("DEBUG: Template data: %+v", templateData)
t.Logf("DEBUG: Rendered template path: %s", templatePath)
```

### Manual Testing

Test VHS templates manually:
```bash
cd features/neovim/{area}/{feature}
# Edit demo.tape.tmpl to use absolute paths
vhs demo.tape.tmpl
```

## Dependencies

- **VHS** (Video Home System) by Charmbracelet
- **notedown-language-server** (built automatically)
- **Neovim** with Lua support
- **Go** for test framework and template processing
- **Git** for golden file management

## Best Practices

1. **Keep tests focused** - One feature per test
2. **Use descriptive names** - Clear area/feature naming
3. **Add proper sleeps** - Allow time for UI updates
4. **Test edge cases** - Include error scenarios
5. **Document thoroughly** - Both code and user documentation
6. **Verify visually** - Check generated GIFs match expectations
7. **Clean workspaces** - Include only necessary test files
8. **Proper timeouts** - Balance speed vs reliability