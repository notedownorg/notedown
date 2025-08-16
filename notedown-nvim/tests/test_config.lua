-- Copyright 2024 Notedown Authors
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

-- Tests for notedown configuration and workspace detection

local MiniTest = require("mini.test")
local utils = require("helpers.utils")

local T = MiniTest.new_set()

T["config defaults"] = function()
	local child = utils.new_child_neovim()

	-- Test that config module loads
	local has_config = child.lua_get('pcall(require, "notedown.config")')
	MiniTest.expect.equality(has_config, true)

	-- Test server defaults
	local server_name = child.lua_get('require("notedown.config").defaults.server.name')
	MiniTest.expect.equality(server_name, "notedown")

	local parser_mode = child.lua_get('require("notedown.config").defaults.parser.mode')
	MiniTest.expect.equality(parser_mode, "auto")

	child.stop()
end

T["workspace detection"] = MiniTest.new_set()

T["workspace detection"]["detects configured workspace"] = function()
	local workspace_path = utils.create_test_workspace("/tmp/test-notedown-workspace")
	local child = utils.new_child_neovim()

	-- Change to the test workspace
	child.lua('vim.fn.chdir("' .. workspace_path .. '")')

	-- Set up notedown with our test workspace
	child.lua(
		'require("notedown").setup({ server = { cmd = { "echo", "mock-server" } }, parser = { mode = "auto", notedown_workspaces = { "'
			.. workspace_path
			.. '" } } })'
	)

	-- Create a markdown buffer
	child.lua('vim.cmd("edit README.md")')

	-- Get workspace status
	local status = child.lua_get('require("notedown").get_workspace_status()')

	MiniTest.expect.equality(status.is_notedown_workspace, true)
	MiniTest.expect.equality(status.should_use_notedown, true)
	-- Check that workspace_path contains our test path (handle macOS /private/tmp vs /tmp)
	local expected_path = workspace_path:gsub("^/tmp", "/private/tmp")
	local actual_path = status.workspace_path or ""
	local path_matches = (actual_path == workspace_path) or (actual_path == expected_path)
	MiniTest.expect.equality(path_matches, true)

	child.stop()
	utils.cleanup_test_workspace(workspace_path)
end

T["workspace detection"]["ignores non-configured workspace"] = function()
	local workspace_path = utils.create_test_workspace("/tmp/test-other-workspace")
	local child = utils.new_child_neovim()

	-- Change to the test workspace
	child.lua('vim.fn.chdir("' .. workspace_path .. '")')

	-- Set up notedown with different workspace
	child.lua(
		'require("notedown").setup({ server = { cmd = { "echo", "mock-server" } }, parser = { mode = "auto", notedown_workspaces = { "/tmp/different-workspace" } } })'
	)

	-- Create a markdown buffer
	child.lua('vim.cmd("edit README.md")')

	-- Get workspace status
	local status = child.lua_get('require("notedown").get_workspace_status()')

	MiniTest.expect.equality(status.is_notedown_workspace, false)
	MiniTest.expect.equality(status.should_use_notedown, false)

	child.stop()
	utils.cleanup_test_workspace(workspace_path)
end

T["parser mode"] = MiniTest.new_set()

T["parser mode"]["respects explicit notedown mode"] = function()
	local child = utils.new_child_neovim()

	-- Set up with explicit notedown mode
	child.lua(
		'require("notedown").setup({ server = { cmd = { "echo", "mock-server" } }, parser = { mode = "notedown" } })'
	)

	-- Create a markdown buffer
	child.lua('vim.cmd("edit test.md")')

	-- Get workspace status
	local status = child.lua_get('require("notedown").get_workspace_status()')

	MiniTest.expect.equality(status.should_use_notedown, true)
	MiniTest.expect.equality(status.parser_mode, "notedown")

	child.stop()
end

T["parser mode"]["respects explicit markdown mode"] = function()
	local child = utils.new_child_neovim()

	-- Set up with explicit markdown mode
	child.lua(
		'require("notedown").setup({ server = { cmd = { "echo", "mock-server" } }, parser = { mode = "markdown" } })'
	)

	-- Create a markdown buffer
	child.lua('vim.cmd("edit test.md")')

	-- Get workspace status
	local status = child.lua_get('require("notedown").get_workspace_status()')

	MiniTest.expect.equality(status.should_use_notedown, false)
	MiniTest.expect.equality(status.parser_mode, "markdown")

	child.stop()
end

T["workspace path normalization"] = function()
	local child = utils.new_child_neovim()

	-- Test tilde expansion and path resolution
	child.lua(
		'require("notedown").setup({ server = { cmd = { "echo", "mock-server" } }, parser = { mode = "auto", notedown_workspaces = { "~/test-workspace", "/tmp/absolute-workspace" } } })'
	)

	-- Check that workspaces were expanded
	local expanded_workspaces = child.lua_get('require("notedown").get_workspace_status().configured_workspaces')

	-- Should have expanded tilde
	local found_home = false
	for _, workspace in ipairs(expanded_workspaces) do
		if workspace:match("^/.*test%-workspace") then
			found_home = true
			break
		end
	end

	MiniTest.expect.equality(found_home, true)

	child.stop()
end

return T
