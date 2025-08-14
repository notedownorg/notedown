local M = {}
local config = require("notedown.config")

function M.setup(opts)
  opts = opts or {}
  
  local final_config = vim.tbl_deep_extend("force", config.defaults, opts)
  
  vim.api.nvim_create_autocmd("FileType", {
    pattern = "markdown",
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