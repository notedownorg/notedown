local M = {}
local config = require("notedown.config")

-- Store the final config for use in other functions
local final_config = {}

-- Check if a file path is within any configured notedown workspace
local function is_notedown_workspace(file_path)
    if not final_config.parser or not final_config.parser.notedown_workspaces then
        return false
    end

    local resolved_file = vim.fn.resolve(file_path)

    for _, workspace in ipairs(final_config.parser.notedown_workspaces) do
        local resolved_workspace = vim.fn.resolve(vim.fn.expand(workspace))

        -- Check if file path starts with workspace path
        if resolved_file:find("^" .. vim.pesc(resolved_workspace)) then
            return true, resolved_workspace
        end
    end

    return false
end

-- Determine if notedown parser should be used for a buffer
local function should_use_notedown_parser(bufnr)
    local file_path = vim.api.nvim_buf_get_name(bufnr)

    -- If no file path, default to markdown
    if file_path == "" then
        return false
    end

    -- Respect explicit user preference
    if final_config.parser.mode == "notedown" then
        return true
    elseif final_config.parser.mode == "markdown" then
        return false
    end

    -- Auto mode: use workspace detection
    if final_config.parser.mode == "auto" then
        return is_notedown_workspace(file_path)
    end

    return false
end

-- Get current workspace status for a buffer
function M.get_workspace_status(bufnr)
    bufnr = bufnr or vim.api.nvim_get_current_buf()
    local file_path = vim.api.nvim_buf_get_name(bufnr)

    if file_path == "" then
        return {
            is_notedown_workspace = false,
            parser_mode = final_config.parser.mode,
            should_use_notedown = false,
        }
    end

    local is_workspace, workspace_path = is_notedown_workspace(file_path)
    local should_use_notedown = should_use_notedown_parser(bufnr)

    return {
        file_path = file_path,
        is_notedown_workspace = is_workspace,
        workspace_path = workspace_path,
        parser_mode = final_config.parser.mode,
        should_use_notedown = should_use_notedown,
        configured_workspaces = final_config.parser.notedown_workspaces,
    }
end

function M.setup(opts)
    opts = opts or {}

    final_config = vim.tbl_deep_extend("force", config.defaults, opts)

    -- Set up parser selection based on workspace detection
    vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
        pattern = "*.md",
        callback = function(args)
            local bufnr = args.buf

            -- Determine which parser to use
            if should_use_notedown_parser(bufnr) then
                -- Set filetype to notedown if available, otherwise markdown
                local has_notedown_parser = pcall(vim.treesitter.language.require_language, "notedown")
                if has_notedown_parser then
                    vim.bo[bufnr].filetype = "notedown"
                else
                    vim.bo[bufnr].filetype = "markdown"
                end
            else
                vim.bo[bufnr].filetype = "markdown"
            end
        end,
        priority = 10, -- Higher priority than default filetype detection
    })

    -- Set up LSP for both markdown and notedown filetypes
    vim.api.nvim_create_autocmd("FileType", {
        pattern = { "markdown", "notedown" },
        callback = function()
            vim.lsp.start({
                name = final_config.server.name,
                cmd = final_config.server.cmd,
                root_dir = final_config.server.root_dir(),
                capabilities = final_config.server.capabilities,
                workspace_folders = {
                    {
                        uri = vim.uri_from_fname(vim.fn.getcwd()),
                        name = vim.fs.basename(vim.fn.getcwd()),
                    }
                }
            })
        end,
    })
end

return M

