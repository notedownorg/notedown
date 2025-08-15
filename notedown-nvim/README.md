# 📝 notedown.nvim

A Neovim plugin for [Notedown Flavored Markdown](https://github.com/notedownorg/notedown) with intelligent LSP integration and workspace-aware parser selection.

<!-- TODO: Add screenshot here -->

## ✨ Features

- 🔗 **Wikilink Support**: Intelligent completion and navigation for `[[wikilinks]]`
- 🏠 **Workspace Detection**: Automatically uses notedown parser for configured workspaces
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
  ft = "markdown",
  opts = {
    -- your configuration comes here
    -- or leave it empty to use the default settings
    -- refer to the configuration section below
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
  parser = {
    mode = "auto", -- "auto" | "notedown" | "markdown"
    notedown_workspaces = {
      -- Add your notedown workspace paths here
      -- Note: This includes any child directories
      -- For example if we set ~/notes I could have ~/notes/workspace1, ~/notes/workspace2, etc
      -- and could open neovim in any of those or their children but maintain workspace separation
      "~/notes",
      "~/github.com/notedownorg/notedown",
    },
  },
})
```

### Parser Modes

- **`"auto"`** (default): Automatically detect notedown workspaces and use appropriate parser
- **`"notedown"`**: Always use notedown parser for all markdown files
- **`"markdown"`**: Always use standard markdown parser

### Workspace Configuration

Configure paths where notedown features should be enabled:

```lua
require("notedown").setup({
  parser = {
    notedown_workspaces = {
      "/Users/username/notes",           -- Your personal notes
      "/Users/username/projects/docs",   -- Project documentation
      "/Users/username/obsidian-vault",  -- Obsidian vault
    }
  }
})
```

**Note**: Workspace paths support:
- Tilde expansion (`~/notes` → `/Users/username/notes`)
- Child directory matching (files in subdirectories are included)
- Symlink resolution for reliable path matching

## 🚀 Usage

### Automatic Features

Once configured, the plugin automatically:
- Detects markdown files in configured workspaces
- Starts the notedown language server
- Provides wikilink completion with `[[`
- Enables go-to-definition for wikilinks

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

### Commands

#### `:NotedownWorkspaceStatus`

Check the workspace status for the current buffer:

```
Notedown Workspace Status:
  File: /Users/username/notes/project/ideas.md
  Parser Mode: auto
  In Notedown Workspace: Yes
  Should Use Notedown Parser: Yes
  Matched Workspace: /Users/username/notes
  Configured Workspaces:
    - /Users/username/notes
    - /Users/username/projects/docs
```

#### `:NotedownReload`

Reload the plugin and restart the LSP server:
- Stops existing LSP clients
- Clears module cache
- Reloads configuration
- Restarts language server

## 🔧 Advanced Configuration

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

1. Ensure file is in a configured workspace
2. Check that LSP server is running: `:LspInfo`
3. Try typing `[[` and wait for completion popup

### Parser Issues

1. Check LSP server status: `:LspInfo`
2. Verify parser mode: `:NotedownWorkspaceStatus`
3. Try forcing parser mode: `parser = { mode = "notedown" }`

## 🤝 Contributing

Contributions are welcome! Please see the [main repository](https://github.com/notedownorg/notedown) for contribution guidelines.

## 📄 License

This project is licensed under the Apache License 2.0. See [LICENSE](../LICENSE) for details.

## 🔗 Related Projects

- [notedown](https://github.com/notedownorg/notedown) - The main Notedown language server
- [Obsidian](https://obsidian.md) - For wikilink inspiration