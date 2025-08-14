local M = {}

M.defaults = {
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
            -- Note: This will includes any child directories
            -- For example if we set ~/notes I could have ~/notes/workspace1, ~/notes/workspace2, etc
            -- and could open neovim in any of three or their children but maintain workspace separation
            "~/notes",
            "~/github.com/notedownorg/notedown",
        },
    },
}

return M

