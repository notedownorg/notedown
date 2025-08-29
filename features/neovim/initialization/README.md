# Initialization Area

This area covers plugin loading, setup, and workspace detection features for the Notedown Neovim plugin.

## Features

### [Workspace Status Command](./workspace-status-command/)
![Workspace Status Demo](./workspace-status-command/demo.gif)

The `:NotedownWorkspaceStatus` command provides information about the current workspace status and LSP connection.

**What it does:**
- Displays current workspace detection status
- Shows LSP server connection information
- Provides diagnostic information for troubleshooting

**Usage:**
```vim
:NotedownWorkspaceStatus
```

**Technical details:**
- Implemented as a Neovim command that queries the LSP client
- Shows workspace root detection results
- Displays server capabilities and connection status

## Area Overview

The initialization area ensures that users can:
1. **Load the plugin** automatically when opening markdown files
2. **Connect to LSP server** without manual configuration
3. **Detect workspaces** correctly based on directory structure
4. **Troubleshoot issues** using status commands

## Running Area Tests

```bash
# Run all initialization tests
go test -run TestFeatures/initialization

# Run specific feature
go test -run TestFeatures/initialization/workspace-status-command
```

## Adding New Features

To add a new initialization feature:

1. Create feature directory: `mkdir initialization/my-new-feature`
2. Add test workspace: `mkdir initialization/my-new-feature/workspace`
3. Create VHS demo: `initialization/my-new-feature/demo.tape.tmpl`
4. Document feature: `initialization/my-new-feature/README.md`
5. Add to test suite in `framework_test.go`

## Related Documentation

- [Plugin Configuration](../../../neovim/README.md)
- [LSP Server Documentation](../../../language-server/)
- [Workspace Detection](../../../pkg/config/)