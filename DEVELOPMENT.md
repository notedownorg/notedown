# Notedown Development Guide

This guide covers how to set up, develop, and test the Notedown project locally.

## Quick Start

### Prerequisites

**Recommended: Nix Package Manager (preferred method)**
- Install [Nix](https://github.com/DeterminateSystems/nix-installer) - provides all dependencies automatically
- The project includes `flake.nix` with complete development environment

**Manual setup (if not using Nix):**
- Go 1.24.4 or later
- Neovim with Lua support
- Git
- VHS (for feature testing)
- golangci-lint (for linting)

### Install Development Build

To install the latest development version (HEAD) on your local machine:

```bash
# Clone the repository
git clone https://github.com/notedownorg/notedown.git
cd notedown

# Install development build
make install
```

This will:
1. **Clean** any existing installation
2. **Build** the LSP server with current version info from git
3. **Install** the binary to `$GOPATH/bin/notedown-language-server`
4. **Install** the Neovim plugin to `~/.config/notedown/nvim/`

### Verify Installation

Check that the installation worked:

```bash
# Verify LSP server binary
which notedown-language-server
notedown-language-server --version

# Verify Neovim plugin files
ls -la ~/.config/notedown/nvim/
```

### Configure Neovim Plugin

After installation, you need to configure Neovim to load the plugin. The method depends on your plugin manager:

#### Lazy.nvim

Add to your Lazy plugin configuration:

```lua
return {
    dir = "~/.config/notedown/nvim",
    name = "notedown",
    config = function()
        vim.treesitter.language.register('markdown', 'notedown')
        require("notedown").setup({})
    end,
}
```

#### Packer.nvim

```lua
use {
    "~/.config/notedown/nvim",
    as = "notedown",
    config = function()
        vim.treesitter.language.register('markdown', 'notedown')
        require("notedown").setup({})
    end,
}
```
### Test with Development Environment

Use the built-in development workspace:

```bash
# Launch development environment
make dev
```

This will:
1. Install the development build
2. Copy the demo workspace to `/tmp/notedown_demo_workspace`
3. Open Neovim in the demo workspace

## Development Workflow

### Making Changes

1. **Make your changes** to the codebase
2. **Test your changes** locally:
   ```bash
   make test-features-fast  # Quick tests without GIF generation
   make test                # Full test suite
   ```
3. **Install and test** your changes:
   ```bash
   make install  # Reinstall with your changes
   make dev      # Test in development environment
   ```

### Code Quality

Before committing changes:

```bash
# Format and tidy code
make hygiene

# Run linter (if available)
make lint

# Run all tests
make all
```

### Version Information

The `make install` command embeds version information from git:

- **Version**: `git describe --tags --always --dirty`
- **Commit**: Current commit hash
- **Date**: Build timestamp

You can see this information with:
```bash
notedown-language-server --version
```

## Build Targets

### Core Development
- `make install` - Build and install development version (HEAD)
- `make clean` - Remove installed binary and plugin files
- `make dev` - Install and open test workspace in Neovim
- `make all` - Full build pipeline with hygiene checks

### Testing
- `make test` - Run all tests (Go, Neovim, and feature tests)
- `make test-features` - Run feature tests with GIF generation
- `make test-features-fast` - Run feature tests without GIFs (faster)
- `make test-features-golden` - Regenerate golden files for feature tests

### Code Quality
- `make format` - Format Go code and Lua files
- `make mod` - Tidy Go modules
- `make hygiene` - Format and mod tidy
- `make lint` - Run golangci-lint (if available)
- `make dirty` - Check if working tree is clean

## Architecture Overview

### Components

1. **Language Server** (`language-server/`) - LSP implementation for Notedown
2. **Parser** (`pkg/parser/`) - Markdown parser with NFM extensions
3. **Neovim Plugin** (`neovim/`) - Lua plugin for Neovim integration
4. **Feature Tests** (`features/neovim/`) - End-to-end testing with VHS

### Installation Details

When you run `make install`:

1. **LSP Server**: Built with ldflags embedding version info, installed to `$GOPATH/bin`
2. **Neovim Plugin**: All files from `neovim/` copied to `~/.config/notedown/nvim/`
3. **Configuration**: Plugin detects the installed LSP server automatically

### File Locations

After installation:
- **LSP Binary**: `$GOPATH/bin/notedown-language-server`
- **Plugin Files**: `~/.config/notedown/nvim/`
  - `lua/notedown/init.lua` - Main plugin code
  - `lua/notedown/config.lua` - Configuration
  - `plugin/notedown.lua` - Plugin bootstrapping

