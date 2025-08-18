# 📝 notedown.nvim

A Neovim plugin for [Notedown Flavored Markdown](https://github.com/notedownorg/notedown) with intelligent LSP integration and workspace-aware parser selection.

<!-- TODO: Add screenshot here -->

## ✨ Features

- 🔗 **Wikilink Support**: Intelligent completion and navigation for `[[wikilinks]]`
- 📝 **List Movement**: Reorganize list items with mk/mj keybindings and smart cursor following
- 🏠 **Automatic Workspace Detection**: Uses notedown parser when `.notedown/` directory is found
- 🧠 **Smart LSP Integration**: Seamless language server integration with document synchronization
- 🚀 **LSP Integration**: Full Notedown Language Server Protocol support
- ⚡ **Fast**: Efficient workspace detection with path-based matching
- 🔧 **Configurable**: Flexible parser selection modes and workspace configuration

## ⚡️ Requirements

- Neovim >= 0.9.0
- [notedown-language-server](https://github.com/notedownorg/notedown) (built and available in PATH)
- TreeSitter with markdown parser installed (for folding support in notedown files)

## 📦 Installation

### [lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
{
  "notedownorg/notedown.nvim",
  opts = {
    -- Most users need no configuration!
    -- Just create a .notedown/ directory in your project root
  }
}

```

### [packer.nvim](https://github.com/wbthomason/packer.nvim)

```lua
use {
  "notedownorg/notedown.nvim",
  config = function()
    require("notedown").setup()
  end
}
```

### [vim-plug](https://github.com/junegunn/vim-plug)

```vim
Plug 'notedownorg/notedown.nvim'

" Then in your init.lua or a lua file
lua require("notedown").setup()
```

## ⚙️ Configuration

### Default Configuration

```lua
require("notedown").setup({
  server = {
    name = "notedown",
    cmd = { "notedown-language-server", "serve", "--log-level", "debug", "--log-file", "/tmp/notedown.log" },
    root_dir = function()
      return vim.fn.getcwd()
    end,
    capabilities = vim.lsp.protocol.make_client_capabilities(),
  },
  keybindings = {
    -- Keybindings for list item movement
    move_list_item_up = "mk",   -- Move up
    move_list_item_down = "mj", -- Move down
  },
})
```

### Workspace Detection

**The plugin automatically detects Notedown workspaces** by looking for a `.notedown/` directory:

1. **Starting from the current file's directory**, the plugin walks up the directory tree
2. **First `.notedown/` directory found** marks the workspace root  
3. **Files in detected workspaces** automatically use the notedown parser and LSP features
4. **Files outside workspaces** use standard markdown behavior

**No configuration required** - just create a `.notedown/` directory in your project root!

#### Setting Up a Workspace

```bash
# Navigate to your project/notes directory
cd ~/my-notes

# Create .notedown directory to mark it as a workspace
mkdir .notedown

# Optional: Add workspace configuration
echo "tasks:" > .notedown/settings.yaml
echo "  states:" >> .notedown/settings.yaml
echo "    - value: ' '" >> .notedown/settings.yaml
echo "      name: todo" >> .notedown/settings.yaml
```

#### Workspace Detection Examples

```
Project Structure:
├── my-notes/
│   ├── .notedown/           ← Workspace root detected here
│   ├── daily/
│   │   └── today.md         ← Opens as notedown
│   ├── projects/
│   │   └── work.md          ← Opens as notedown  
│   └── README.md            ← Opens as notedown

├── other-project/
│   └── README.md            ← Opens as markdown (no .notedown)
```

## 🚀 Usage

### Automatic Features

The plugin automatically:
- Detects Notedown workspaces by finding `.notedown/` directories
- Starts the notedown language server for markdown files in workspaces
- Provides wikilink completion with `[[`
- Enables go-to-definition for wikilinks
- Sets the LSP root directory to the detected workspace root

### LSP Features

#### Wikilink Completion

Type `[[` to trigger intelligent completion:

- **Existing Files**: Complete paths to actual markdown files
- **Referenced Targets**: Suggest wikilink targets mentioned in other files
- **Directory Paths**: Complete directory structures for organization

#### Go-to-Definition

- Place cursor on a wikilink target
- Use `gd` or your configured go-to-definition keybinding
- Jump to the target file or create it if it doesn't exist

### List Movement

Reorganize list items quickly with intuitive keybindings:

- **`mk`**: Move list item up  
- **`mj`**: Move list item down

Features:
- **Smart cursor following**: Cursor stays with the moved content
- **Multi-level support**: Works with nested lists of any depth
- **List type aware**: Handles bullet lists, numbered lists, and task lists
- **Auto-renumbering**: Updates list numbers when moving numbered items
- **Character position preservation**: Maintains cursor position within moved text

Example:
```markdown
1. First item
2. Second item   <- cursor here, press mk
3. Third item
```

Becomes:
```markdown
1. First item
2. Third item
3. Second item   <- cursor follows the moved item
```

### Commands

#### `:NotedownWorkspaceStatus`

Check the workspace status for the current buffer:

```
Notedown Workspace Status:
  File: /Users/username/notes/ideas.md
  In Notedown Workspace: Yes
  Should Use Notedown Parser: Yes
  Matched Workspace: /Users/username/notes
  Detection Method: Auto-detected (.notedown directory)
```

#### `:NotedownReload`

Reload the plugin and restart the LSP server:
- Stops existing LSP clients
- Clears module cache
- Reloads configuration
- Restarts language server

## 🔧 Advanced Configuration

### Custom Keybindings

Customize list movement keybindings to your preference:

```lua
require("notedown").setup({
  keybindings = {
    move_list_item_up = "[e",     -- Use bracket notation
    move_list_item_down = "]e",   -- Use bracket notation
    -- or
    move_list_item_up = "<leader>k",   -- Use leader-based
    move_list_item_down = "<leader>j", -- Use leader-based
    -- or
    move_list_item_up = "<M-k>",   -- Use Alt (may need terminal config)
    move_list_item_down = "<M-j>", -- Use Alt (may need terminal config)
  },
})
```

### Custom LSP Server Command

```lua
require("notedown").setup({
  server = {
    cmd = { "/path/to/notedown-language-server", "serve", "--log-level", "info" },
    root_dir = function()
      -- Use git root or fallback to current directory
      return vim.fn.system("git rev-parse --show-toplevel"):gsub("\n", "") or vim.fn.getcwd()
    end,
  }
})
```


### Custom Capabilities

```lua
require("notedown").setup({
  server = {
    capabilities = vim.tbl_deep_extend(
      "force",
      vim.lsp.protocol.make_client_capabilities(),
      require("cmp_nvim_lsp").default_capabilities() -- if using nvim-cmp
    ),
  }
})
```

## 🐛 Troubleshooting

### LSP Server Not Starting

1. Ensure `notedown-language-server` is in your PATH:
   ```bash
   which notedown-language-server
   ```

2. Check server logs:
   ```bash
   tail -f /tmp/notedown.log
   ```

3. Verify configuration with `:NotedownWorkspaceStatus`

### Wikilink Completion Not Working

1. Ensure a `.notedown/` directory exists in your project root or parent directories
2. Check that LSP server is running: `:LspInfo`
3. Verify workspace detection: `:NotedownWorkspaceStatus`
4. Try typing `[[` and wait for completion popup

### Parser Issues

1. Check LSP server status: `:LspInfo`
2. Verify workspace detection: `:NotedownWorkspaceStatus`
3. Ensure `.notedown/` directory exists in your project root or a parent directory

## 🤝 Contributing

Contributions are welcome! Please see the [main repository](https://github.com/notedownorg/notedown) for contribution guidelines.

## 📄 License

This project is licensed under the Apache License 2.0. See [LICENSE](../LICENSE) for details.

## 🔗 Related Projects

- [notedown](https://github.com/notedownorg/notedown) - The main Notedown language server
- [Obsidian](https://obsidian.md) - For wikilink inspiration
