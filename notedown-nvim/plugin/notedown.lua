if vim.g.loaded_notedown then
  return
end
vim.g.loaded_notedown = 1

vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
  pattern = "*.md",
  callback = function()
    vim.bo.filetype = "markdown"
  end,
})

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