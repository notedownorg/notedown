-- init.lua for VHS testing - copied to ~/.config/nvim/init.lua in Docker
-- Set up basic Neovim configuration (modern syntax for Neovim 0.7+)
vim.opt.number = true
vim.opt.relativenumber = false
vim.opt.wrap = false
vim.opt.termguicolors = true  -- Enable for colorscheme

-- Add the built-in neovim plugin directory to runtime path
vim.opt.runtimepath:prepend('/opt/notedown/nvim')

-- Configure Catppuccin theme (installed via git clone in Docker)
require("catppuccin").setup({
    flavour = "mocha", -- Use the dark mocha variant
    background = {
        dark = "mocha",
    },
    transparent_background = false,
    integrations = {
        treesitter = true,
        native_lsp = {
            enabled = true,
        },
    },
})

-- Set the colorscheme
vim.cmd.colorscheme("catppuccin-mocha")

-- Use native Neovim completion for VHS demonstrations
-- Configure built-in LSP completion
vim.opt.completeopt = { 'menu', 'menuone', 'noselect' }

-- Configure Notedown plugin with container paths
require('notedown').setup({
    server = {
        cmd = {'/usr/local/bin/notedown-language-server', 'serve'},
        auto_start = true,
    },
    workspaces = {
        {
            pattern = '/vhs/workspace',
            parser = 'notedown'
        }
    }
})

-- Let Notedown plugin handle concealment settings automatically
-- The plugin sets conceallevel=2 and concealcursor="nc" for notedown files