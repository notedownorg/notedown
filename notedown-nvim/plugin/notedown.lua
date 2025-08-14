if vim.g.loaded_notedown then
  return
end
vim.g.loaded_notedown = 1


-- Note: Filetype detection is now handled in init.lua based on workspace detection
-- This autocmd is kept for compatibility but may be overridden

vim.api.nvim_create_user_command("NotedownReload", function()
  -- Stop existing LSP clients
  vim.lsp.stop_client(vim.lsp.get_active_clients({ name = "notedown" }))
  
  -- Clear module cache
  package.loaded['notedown'] = nil
  package.loaded['notedown.config'] = nil
  package.loaded['notedown.init'] = nil
  
  -- Reload the plugin
  require('notedown').setup()
  
  -- If current buffer is markdown, trigger the autocmd
  if vim.bo.filetype == "markdown" then
    vim.api.nvim_exec_autocmds("FileType", { pattern = "markdown" })
  end
  
  vim.notify("Notedown plugin and LSP reloaded", vim.log.levels.INFO)
end, {
  desc = "Reload the Notedown plugin and restart LSP",
})

vim.api.nvim_create_user_command("NotedownWorkspaceStatus", function()
  local notedown = require('notedown')
  local status = notedown.get_workspace_status()
  
  local message = string.format([[
Notedown Workspace Status:
  File: %s
  Parser Mode: %s
  In Notedown Workspace: %s
  Should Use Notedown Parser: %s
  ]], 
    status.file_path or "No file",
    status.parser_mode,
    status.is_notedown_workspace and "Yes" or "No",
    status.should_use_notedown and "Yes" or "No"
  )
  
  if status.workspace_path then
    message = message .. string.format("  Matched Workspace: %s\n", status.workspace_path)
  end
  
  if status.configured_workspaces and #status.configured_workspaces > 0 then
    message = message .. "  Configured Workspaces:\n"
    for _, workspace in ipairs(status.configured_workspaces) do
      message = message .. string.format("    - %s\n", workspace)
    end
  else
    message = message .. "  No workspaces configured\n"
  end
  
  print(message)
end, {
  desc = "Show workspace status for current buffer",
})