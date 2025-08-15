local M = {}
local config = require("notedown.config")

-- Store the final config for use in other functions
local final_config = {}

-- Check if current working directory matches any configured notedown workspace
local function is_notedown_workspace(file_path)
    if not final_config.parser or not final_config.parser.notedown_workspaces then
        return false
    end

    local cwd = vim.fn.resolve(vim.fn.getcwd())

    for _, workspace in ipairs(final_config.parser.notedown_workspaces) do
        -- Workspace paths are already expanded and resolved during setup
        -- Check if current working directory exactly matches workspace path
        if cwd == workspace then
            return true, workspace
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

    local is_workspace, workspace_path = is_notedown_workspace()
    local should_use_notedown = should_use_notedown_parser(bufnr)

    return {
        file_path = file_path,
        cwd = vim.fn.getcwd(),
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

    -- Expand and normalize workspace paths during setup
    if final_config.parser and final_config.parser.notedown_workspaces then
        local expanded_workspaces = {}
        for _, workspace in ipairs(final_config.parser.notedown_workspaces) do
            -- Expand ~ and resolve the path
            local expanded = vim.fn.expand(workspace)
            local resolved = vim.fn.resolve(expanded)
            table.insert(expanded_workspaces, resolved)
        end
        
        -- Filter out child directories and notify about ignored paths
        local filtered_workspaces = {}
        local ignored_paths = {}
        
        for _, workspace in ipairs(expanded_workspaces) do
            local is_child = false
            
            -- Check if this workspace is a child of any existing workspace
            for _, existing in ipairs(filtered_workspaces) do
                if workspace:find("^" .. vim.pesc(existing) .. "/") then
                    is_child = true
                    table.insert(ignored_paths, workspace)
                    break
                end
            end
            
            if not is_child then
                -- Check if any existing workspaces are children of this one
                local children_to_remove = {}
                for i, existing in ipairs(filtered_workspaces) do
                    if existing:find("^" .. vim.pesc(workspace) .. "/") then
                        table.insert(children_to_remove, i)
                        table.insert(ignored_paths, existing)
                    end
                end
                
                -- Remove children (in reverse order to maintain indices)
                for i = #children_to_remove, 1, -1 do
                    table.remove(filtered_workspaces, children_to_remove[i])
                end
                
                table.insert(filtered_workspaces, workspace)
            end
        end
        
        -- Notify about ignored paths
        if #ignored_paths > 0 then
            vim.notify(
                "Notedown: Ignored workspace paths (child directories of other workspaces):\n  " .. 
                table.concat(ignored_paths, "\n  "),
                vim.log.levels.WARN
            )
        end
        
        final_config.parser.notedown_workspaces = filtered_workspaces
    end

    -- Set up parser selection based on workspace detection
    vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
        pattern = "*.md",
        callback = function(args)
            local bufnr = args.buf

            -- Determine which filetype to use
            if should_use_notedown_parser(bufnr) then
                vim.bo[bufnr].filetype = "notedown"
            else
                vim.bo[bufnr].filetype = "markdown"
            end
        end,
    })

    -- Set up LSP and folding for both markdown and notedown filetypes
    vim.api.nvim_create_autocmd("FileType", {
        pattern = { "markdown", "notedown" },
        callback = function()
            -- Start LSP
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

            -- Enable treesitter-based folding for notedown files
            if vim.bo.filetype == "notedown" then
                vim.opt_local.foldmethod = "expr"
                vim.opt_local.foldexpr = "v:lua.vim.treesitter.foldexpr()"
                vim.opt_local.foldenable = true
                vim.opt_local.foldlevel = 99 -- Start with all folds open
            end

            -- Set up keybindings for list item movement (for both markdown and notedown)
            M.setup_list_movement_keybindings()
        end,
    })
end

-- Get the appropriate notedown LSP client for command execution
local function get_notedown_command_client()
    local clients = vim.lsp.get_active_clients({ name = "notedown" })
    if #clients == 0 then
        vim.notify("Notedown LSP server not active", vim.log.levels.WARN)
        return nil
    end
    
    -- Find the client that supports executeCommand
    for _, client in ipairs(clients) do
        if client.server_capabilities and client.server_capabilities.executeCommandProvider then
            return client
        end
    end
    
    vim.notify("No notedown client supports executeCommand", vim.log.levels.WARN)
    return nil
end

-- Set up keybindings for list item movement
function M.setup_list_movement_keybindings()
    local opts = { buffer = true, silent = true }

    -- Move list item up
    vim.keymap.set('n', final_config.keybindings.move_list_item_up, function()
        M.move_list_item_up()
    end, vim.tbl_extend('force', opts, { desc = 'Move list item up' }))

    -- Move list item down
    vim.keymap.set('n', final_config.keybindings.move_list_item_down, function()
        M.move_list_item_down()
    end, vim.tbl_extend('force', opts, { desc = 'Move list item down' }))
end

-- Move list item up via LSP command
function M.move_list_item_up()
    local client = get_notedown_command_client()
    if not client then
        return
    end
    
    local cursor = vim.api.nvim_win_get_cursor(0)
    local position = {
        line = cursor[1] - 1, -- Convert to 0-based
        character = cursor[2]
    }
    
    local params = {
        command = "notedown.moveListItemUp",
        arguments = {
            vim.uri_from_bufnr(0),
            position
        }
    }
    
    vim.notify(string.format("Moving list item up at line %d", cursor[1]), vim.log.levels.DEBUG)
    
    -- Execute command on the specific notedown client
    client.request("workspace/executeCommand", params, function(err, result)
        if err then
            local error_msg = tostring(err)
            if error_msg:match("boundary") then
                vim.notify("Cannot move up: already at top", vim.log.levels.WARN)
            else
                vim.notify("Error moving list item: " .. error_msg, vim.log.levels.ERROR)
            end
            return
        end
        
        if result then
            -- Check if result is a workspace edit
            if result.changes then
                -- Apply the workspace edit
                local success = vim.lsp.util.apply_workspace_edit(result, "utf-8")
                if success then
                    vim.notify("List item moved up", vim.log.levels.INFO)
                else
                    vim.notify("Failed to apply workspace edit", vim.log.levels.ERROR)
                end
            else
                vim.notify("Command completed but no workspace edit returned", vim.log.levels.WARN)
            end
        else
            vim.notify("No move performed", vim.log.levels.INFO)
        end
    end)
end

-- Move list item down via LSP command
function M.move_list_item_down()
    local client = get_notedown_command_client()
    if not client then
        return
    end
    
    local cursor = vim.api.nvim_win_get_cursor(0)
    local position = {
        line = cursor[1] - 1, -- Convert to 0-based
        character = cursor[2]
    }
    
    local params = {
        command = "notedown.moveListItemDown",
        arguments = {
            vim.uri_from_bufnr(0),
            position
        }
    }
    
    vim.notify(string.format("Moving list item down at line %d", cursor[1]), vim.log.levels.DEBUG)
    
    -- Execute command on the specific notedown client
    client.request("workspace/executeCommand", params, function(err, result)
        if err then
            local error_msg = tostring(err)
            if error_msg:match("boundary") then
                vim.notify("Cannot move down: already at bottom", vim.log.levels.WARN)
            else
                vim.notify("Error moving list item: " .. error_msg, vim.log.levels.ERROR)
            end
            return
        end
        
        if result then
            -- Check if result is a workspace edit
            if result.changes then
                -- Apply the workspace edit
                local success = vim.lsp.util.apply_workspace_edit(result, "utf-8")
                if success then
                    vim.notify("List item moved down", vim.log.levels.INFO)
                else
                    vim.notify("Failed to apply workspace edit", vim.log.levels.ERROR)
                end
            else
                vim.notify("Command completed but no workspace edit returned", vim.log.levels.WARN)
            end
        else
            vim.notify("No move performed", vim.log.levels.INFO)
        end
    end)
end

return M
