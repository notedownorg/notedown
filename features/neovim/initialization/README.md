# Plugin Initialization

This area covers how the Notedown Neovim plugin starts up, detects your workspace, and connects to the language server.

## Available Features

### [Workspace Status Command](./workspace-status-command/)
![Workspace Status Demo](./workspace-status-command/demo.gif)

The `:NotedownWorkspaceStatus` command helps you understand how Notedown is configured in your current workspace.

**What it shows:**
- Whether you're in a Notedown workspace (has a `.notedown` directory)
- Which parser Notedown will use for your files
- Language server connection status
- How the workspace was detected

**How to use:**
```vim
:NotedownWorkspaceStatus
```

This is particularly useful when troubleshooting why features aren't working as expected.

## What Initialization Covers

When you open a Markdown file, the Notedown plugin:

1. **Detects your workspace** by looking for a `.notedown` directory
2. **Chooses the right parser** (Notedown or standard Markdown)
3. **Connects to the language server** for advanced features
4. **Provides status commands** to help you troubleshoot

## Getting Started

1. Make sure you have the Notedown plugin installed
2. Create or open a directory with a `.notedown` folder (this makes it a Notedown workspace)
3. Open any `.md` file in that workspace
4. Try running `:NotedownWorkspaceStatus` to see the plugin in action

## Troubleshooting

If the plugin isn't working as expected:

- Run `:NotedownWorkspaceStatus` to see current status
- Check if you're in a Notedown workspace (should show "Yes")
- Verify the LSP server is running (should show "Active")
- Look at the detection method to understand how workspace was found