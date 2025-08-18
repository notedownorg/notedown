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
	keybindings = {
		-- Keybindings for list item movement
		move_list_item_up = "mk",
		move_list_item_down = "mj",
	},
}

return M
