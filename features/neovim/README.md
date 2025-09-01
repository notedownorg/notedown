# Notedown Neovim Plugin Features

This directory contains **living documentation** for Notedown's Neovim plugin features. Each feature is demonstrated through real terminal interactions that you can see in action.

## What is Living Documentation?

Unlike static documentation that can become outdated, our living documentation:
- **Shows** features working through visual demonstrations
- **Proves** functionality through automated testing
- **Stays current** by being automatically verified with each code change

## Available Features

### [Initialization](./initialization/)
Plugin loading, setup, and workspace detection features.

- **[Workspace Status Command](./initialization/workspace-status-command/)** - Check your current workspace status with `:NotedownWorkspaceStatus`

### [Wikilinks](./wikilinks/)
Internal linking features for connecting documents within your workspace.

- **[Completion](./wikilinks/completion/)** - Auto-complete wikilink targets as you type
- **[Syntax](./wikilinks/syntax/)** - Comprehensive wikilink syntax including display text, path resolution, navigation, and concealment
- **[Diagnostics and Code Actions](./wikilinks/diagnostics-and-code-actions/)** - Real-time feedback on wikilink validity and automated fixes

## How to Use These Features

Each feature directory contains:
- **README.md** - Detailed explanation of the feature
- **demo.gif** - Visual demonstration showing the feature in action
- Working examples you can try in your own Notedown workspace

### Example: Workspace Status Command

The workspace status command helps you understand how Notedown detects your workspace:

![Workspace Status Demo](./initialization/workspace-status-command/demo.gif)

This command shows:
- Whether you're in a Notedown workspace
- Which parser Notedown will use
- LSP server connection status
- How the workspace was detected

## Getting Started

1. **Install the Notedown Neovim plugin** (see main project README)
2. **Open a Notedown workspace** (any directory with a `.notedown` folder)
3. **Try the features** documented in each area

## Feature Areas

Features are organized by functional area:

- **Initialization** - Plugin setup and workspace detection
- **Wikilinks** - Internal linking and document navigation

Each area contains multiple related features with comprehensive documentation and demonstrations.