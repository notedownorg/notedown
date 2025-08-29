# Workspace Status Command

![Demo](demo.gif)

The `:NotedownWorkspaceStatus` command helps you understand how Notedown sees your current workspace and whether everything is set up correctly.

## What This Command Does

When you run `:NotedownWorkspaceStatus`, it shows you:

- **Are you in a Notedown workspace?** - Whether the current directory has a `.notedown` folder
- **Is the language server running?** - Connection status and number of active clients
- **Which parser will be used?** - Notedown parser vs. standard Markdown parser
- **How was the workspace detected?** - Auto-detection method used

## How to Use It

1. Open any Markdown file in your workspace
2. Run the command:
   ```vim
   :NotedownWorkspaceStatus
   ```
3. Read the status information displayed

## Example Output

```
Notedown Workspace Status:
  File: /path/to/your/workspace/document.md
  In Notedown Workspace: Yes
  Should Use Notedown Parser: Yes
  LSP Server Status: Active (1 clients)
    Matched Workspace: /path/to/your/workspace
  Detection Method: Auto-detected (.notedown directory)

Press ENTER or type command to continue
```

## When to Use This Command

### Getting Started
- **After installing** the Notedown plugin to verify it's working
- **When opening** a new workspace to confirm Notedown detects it correctly
- **If features aren't working** as expected

### Troubleshooting
- **Plugin seems inactive** - Check if you're in a Notedown workspace
- **Features not available** - Verify the language server is connected
- **Wrong parser behavior** - See which parser Notedown is using

## Understanding the Output

### Workspace Detection
- **"Yes"** - You're in a directory with a `.notedown` folder (Notedown features active)
- **"No"** - Standard Markdown directory (limited Notedown features)

### Parser Selection  
- **"Yes"** - Files will be parsed with Notedown-specific features
- **"No"** - Files will be parsed as standard Markdown

### LSP Server Status
- **"Active"** - Language server is running and connected
- **"Inactive"** - Language server issues (check installation)

### Detection Method
- **"Auto-detected"** - Found `.notedown` directory automatically
- **"No .notedown directory found"** - Using standard Markdown mode

## Quick Fixes

**If workspace shows "No":**
1. Create a `.notedown` directory in your project root
2. Restart Neovim or run `:NotedownReload`

**If LSP shows "Inactive":**
1. Check that `notedown-language-server` is installed
2. Try `:NotedownReload` to restart the plugin

**If wrong parser is selected:**
1. Verify you're in the correct workspace directory
2. Check that `.notedown` directory exists and is accessible