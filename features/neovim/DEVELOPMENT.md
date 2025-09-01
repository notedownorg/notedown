# Feature Development Guide

This guide covers how to develop, test, and maintain the Notedown Neovim plugin feature documentation system.

## Architecture

### Directory Structure

```
features/neovim/
├── README.md                           # User-focused feature overview
├── DEVELOPMENT.md                      # This file - development guide
├── framework_test.go                   # Test runner for all features
├── pkg/notedown/container_runner.go    # Containerized testing framework
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

The `ContainerVHSRunner` provides containerized testing with Docker:

- **Docker Image Building**: Shared build with `sync.Once` coordination
- **Container Isolation**: Each test runs in a clean Docker environment  
- **Volume Mounting**: Test workspaces mounted into containers
- **Template Rendering**: Go template-based VHS tape generation
- **Parallel Execution**: Concurrent container runs with proper isolation

## Running Tests

### Prerequisites
```bash
# Build the Docker image first (required for containerized testing)
docker build -t notedown-vhs:latest -f features/neovim/Dockerfile .
```

### All Features
```bash
# Run all feature tests (containerized)
make test-features

# Run specific area/feature
go test -v -run TestFeatures/initialization/workspace-status-command

# Run with GIF generation disabled
go test -gif=false -v -run TestFeatures/initialization/workspace-status-command
```

### Development Workflow
```bash
# Build Docker image (if not done already)
docker build -t notedown-vhs:latest -f features/neovim/Dockerfile .

# Run a single test during development 
cd features/neovim
go test -v -run TestFeatures/initialization/workspace-status-command

# Run single test without GIF generation (faster)
go test -gif=false -v -run TestFeatures/initialization/workspace-status-command

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
- GIF generation is configurable via `-gif` flag (default: true)
- **With GIFs** (`-gif=true`): Full documentation generation with visual demonstrations
- **Without GIFs** (`-gif=false`): Fast testing mode - ASCII only (much faster)
- Both ASCII and GIF outputs generated simultaneously when enabled
- **Note**: GIF generation happens inside Docker containers with virtual display

### Parallel Execution
- Tests run concurrently for faster feedback
- Shared Docker image building with `sync.Once` coordination
- Per-test isolated containers with volume mounts
- Container cleanup after each test

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

# Start Neovim (plugin auto-loaded via ~/.config/nvim/init.lua)
Type "nvim ."
Enter
Sleep 4s

# Test feature
Type ":NotedownWorkspaceStatus"
Enter
Sleep 3s
```

### Available Template Data
- `OutputFile` - Container path for ASCII output (`/vhs/{test}.ascii`)
- `WorkspaceDir` - Container workspace directory (`/vhs/workspace`)
- `TmpDir` - Container temporary directory (`/vhs`)
- `LSPBinary` - Container LSP binary path (`/usr/local/bin/notedown-language-server`)
- `TestName` - Safe test name for file operations
- `GenerateGIF` - Boolean flag for GIF generation

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

# Start Neovim (plugin auto-loads)
Type "nvim index.md"
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
- Build the Docker image (if needed)
- Create a container with your workspace
- Run your VHS template in the container
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

### ContainerVHSRunner

The containerized test runner handles:

1. **Docker Image Building**: Shared build using `sync.Once` coordination
2. **Container Management**: Creates and manages test containers with testcontainers-go
3. **Volume Mounting**: Mounts test workspaces into containers
4. **Template Rendering**: Processes Go templates with container paths
5. **VHS Execution**: Runs VHS inside containers with virtual displays
6. **Asset Management**: Copies GIFs and manages golden files

### Container Environment

The Docker container provides:
1. **Debian bookworm-slim** base with essential tools
2. **Neovim 0.10.2** installed from GitHub releases
3. **VHS** and dependencies for terminal recording
4. **LSP server binary** built and installed
5. **Neovim plugin** pre-installed in `/opt/notedown/nvim/`
6. **Configuration** auto-loaded from `~/.config/nvim/init.lua`

### Template Processing

Templates are processed with:
1. Safe filename generation (replacing `/` with `-`)
2. Container path resolution for mounted volumes
3. Dynamic data injection for container paths
4. Error handling for missing templates

## CI Integration

### GitHub Actions

Tests run in the `features` job in `.github/workflows/test.yml`:
- Uses Docker for containerized VHS testing
- Builds notedown-vhs image once per CI run
- Runs in parallel with other test suites
- Generates ASCII output for regression testing
- GIFs generated in headless containers with virtual displays

### Quality Checks

Before committing:
```bash
make format  # Code formatting
make mod     # Go module tidying
```

## Troubleshooting

### Common Issues

1. **Docker image not found**: Run `docker build -t notedown-vhs:latest .` first
2. **Container startup failures**: Check Docker daemon is running
3. **Template errors**: Verify all template variables use container paths
4. **Golden file mismatches**: Delete and regenerate golden files
5. **VHS timeouts**: Increase timeout in test definition
6. **Mount issues**: Ensure workspace directories exist before testing

### Debug Mode

Add debug output to runner:
```go
t.Logf("DEBUG: Template data: %+v", templateData)
t.Logf("DEBUG: Rendered template path: %s", templatePath)
```

### Manual Testing

Test VHS templates manually with Docker:
```bash
# Build image (from project root)
docker build -t notedown-vhs:latest -f features/neovim/Dockerfile .

# Run specific test manually
cd features/neovim/{area}/{feature}
docker run --rm -v "$PWD:/vhs" notedown-vhs:latest demo.tape.tmpl
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