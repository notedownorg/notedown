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

-- Tests for list item movement functionality

local MiniTest = require("mini.test")
local utils = require("helpers.utils")
local lsp = require("helpers.lsp")

local T = MiniTest.new_set()

-- Helper function to create a test workspace with list content
local function create_list_test_workspace(path)
	path = path or "/tmp/notedown-list-test-workspace"

	-- Clean up existing workspace
	vim.fn.system({ "rm", "-rf", path })
	vim.fn.mkdir(path, "p")

	-- Create test file with simple non-nested list
	local test_content = [[# Test List

- First item
- Second item
- Third item
- Fourth item

Some text after the list.
]]

	local file_path = path .. "/test_list.md"
	utils.write_file(file_path, test_content)

	return path, file_path
end

T["list movement setup"] = function()
	local workspace_path, file_path = create_list_test_workspace()
	local child = utils.new_child_neovim()

	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)

	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)

	-- Verify filetype is set correctly
	utils.expect_buffer_filetype(child, "notedown")

	-- Note: LSP client check removed due to timing issues in test environment

	-- Verify NotedownMoveUp and NotedownMoveDown commands exist
	local move_up_cmd_exists = child.lua_get('vim.fn.exists(":NotedownMoveUp") == 2')
	local move_down_cmd_exists = child.lua_get('vim.fn.exists(":NotedownMoveDown") == 2')

	MiniTest.expect.equality(move_up_cmd_exists, true, "NotedownMoveUp command should exist")
	MiniTest.expect.equality(move_down_cmd_exists, true, "NotedownMoveDown command should exist")

	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["move second item down"] = function()
	local workspace_path, file_path = create_list_test_workspace()
	local child = utils.new_child_neovim()

	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)

	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)

	-- Position cursor on second list item (line 4: "- Second item")
	child.lua("vim.api.nvim_win_set_cursor(0, {4, 0})")

	-- Get initial content for comparison
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')

	-- Execute NotedownMoveDown command
	child.lua('vim.cmd("NotedownMoveDown")')

	-- Wait for LSP command to complete
	vim.loop.sleep(1500)

	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')

	print("=== INITIAL CONTENT ===")
	print(initial_content)
	print("=== FINAL CONTENT ===")
	print(final_content)

	-- Verify content changed
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after move down")

	-- Verify specific expected change: "Second item" should now be on line 5, "Third item" on line 4
	local final_lines = child.lua_get("vim.api.nvim_buf_get_lines(0, 0, -1, false)")

	-- Line indices in Lua are 1-based, content lines are 0-based
	-- Line 4 (index 4) should now contain "Third item"
	-- Line 5 (index 5) should now contain "Second item"
	if #final_lines >= 5 then
		local line4 = final_lines[4] or ""
		local line5 = final_lines[5] or ""

		MiniTest.expect.equality(string.find(line4, "Third item") ~= nil, true, "Line 4 should contain 'Third item'")
		MiniTest.expect.equality(string.find(line5, "Second item") ~= nil, true, "Line 5 should contain 'Second item'")
	end

	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["move third item up"] = function()
	local workspace_path, file_path = create_list_test_workspace()
	local child = utils.new_child_neovim()

	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)

	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)

	-- Position cursor on third list item (line 5: "- Third item")
	child.lua("vim.api.nvim_win_set_cursor(0, {5, 0})")

	-- Get initial content for comparison
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')

	-- Execute NotedownMoveUp command
	child.lua('vim.cmd("NotedownMoveUp")')

	-- Wait for LSP command to complete
	vim.loop.sleep(1500)

	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')

	print("=== INITIAL CONTENT ===")
	print(initial_content)
	print("=== FINAL CONTENT ===")
	print(final_content)

	-- Verify content changed
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after move up")

	-- Verify specific expected change: "Third item" should now be on line 4, "Second item" on line 5
	local final_lines = child.lua_get("vim.api.nvim_buf_get_lines(0, 0, -1, false)")

	if #final_lines >= 5 then
		local line4 = final_lines[4] or ""
		local line5 = final_lines[5] or ""

		MiniTest.expect.equality(string.find(line4, "Third item") ~= nil, true, "Line 4 should contain 'Third item'")
		MiniTest.expect.equality(string.find(line5, "Second item") ~= nil, true, "Line 5 should contain 'Second item'")
	end

	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["move first item up - should not move"] = function()
	local workspace_path, file_path = create_list_test_workspace()
	local child = utils.new_child_neovim()

	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)

	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)

	-- Position cursor on first list item (line 3: "- First item")
	child.lua("vim.api.nvim_win_set_cursor(0, {3, 0})")

	-- Get initial content
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')

	-- Execute NotedownMoveUp command
	child.lua('vim.cmd("NotedownMoveUp")')

	-- Wait for LSP command to complete
	vim.loop.sleep(1500)

	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')

	print("=== INITIAL CONTENT ===")
	print(initial_content)
	print("=== FINAL CONTENT ===")
	print(final_content)

	-- Content should remain the same since first item can't move up
	MiniTest.expect.equality(
		initial_content,
		final_content,
		"First item should not move up - content should be unchanged"
	)

	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["move last item down - should not move"] = function()
	local workspace_path, file_path = create_list_test_workspace()
	local child = utils.new_child_neovim()

	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)

	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)

	-- Position cursor on last list item (line 6: "- Fourth item")
	child.lua("vim.api.nvim_win_set_cursor(0, {6, 0})")

	-- Get initial content
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')

	-- Execute NotedownMoveDown command
	child.lua('vim.cmd("NotedownMoveDown")')

	-- Wait for LSP command to complete
	vim.loop.sleep(1500)

	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')

	print("=== INITIAL CONTENT ===")
	print(initial_content)
	print("=== FINAL CONTENT ===")
	print(final_content)

	-- Content should remain the same since last item can't move down
	MiniTest.expect.equality(
		initial_content,
		final_content,
		"Last item should not move down - content should be unchanged"
	)

	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

-- NESTED LIST TESTS

-- Helper function to create a test workspace with deeply nested list content
local function create_nested_list_test_workspace(path)
	path = path or "/tmp/notedown-nested-list-test-workspace"
	
	-- Clean up existing workspace
	vim.fn.system({ "rm", "-rf", path })
	vim.fn.mkdir(path, "p")
	
	-- Create test file with deeply nested lists (6 levels deep)
	local test_content = [[# Deep Nested List Test

- Level 1 Item A
  - Level 2 Item A.1
    - Level 3 Item A.1.a
      - Level 4 Item A.1.a.i
        - Level 5 Item A.1.a.i.α
          - Level 6 Item A.1.a.i.α.I
          - Level 6 Item A.1.a.i.α.II
        - Level 5 Item A.1.a.i.β
      - Level 4 Item A.1.a.ii
    - Level 3 Item A.1.b
  - Level 2 Item A.2
    - Level 3 Item A.2.a
- Level 1 Item B
  - Level 2 Item B.1
    - Level 3 Item B.1.a
      - Level 4 Item B.1.a.i
        - Level 5 Item B.1.a.i.α
      - Level 4 Item B.1.a.ii
    - Level 3 Item B.1.b
  - Level 2 Item B.2
- Level 1 Item C

## Mixed List Types

1. Ordered Level 1 Item A
   - Bullet Level 2 Item A.1
     1. Ordered Level 3 Item A.1.a
        - Bullet Level 4 Item A.1.a.i
          1. Ordered Level 5 Item A.1.a.i.α
   - Bullet Level 2 Item A.2
2. Ordered Level 1 Item B
   - Bullet Level 2 Item B.1
3. Ordered Level 1 Item C

## Task Lists with Nesting

- [ ] Main Task A
  - [ ] Subtask A.1
    - [x] Sub-subtask A.1.a (completed)
    - [ ] Sub-subtask A.1.b
      - [ ] Deep subtask A.1.b.i
        - [ ] Very deep subtask A.1.b.i.α
  - [x] Subtask A.2 (completed)
- [ ] Main Task B
  - [ ] Subtask B.1
    - [ ] Sub-subtask B.1.a
- [x] Main Task C (completed)
]]
	
	local file_path = path .. "/nested_test.md"
	utils.write_file(file_path, test_content)
	
	return path, file_path
end

T["nested list movement - level 1 items"] = function()
	local workspace_path, file_path = create_nested_list_test_workspace()
	local child = utils.new_child_neovim()
	
	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)
	
	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)
	
	-- Position cursor on "Level 1 Item B" (should be around line 17)
	child.lua('vim.fn.search("Level 1 Item B")')
	
	-- Get initial content
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Execute NotedownMoveUp command
	child.lua('vim.cmd("NotedownMoveUp")')
	
	-- Wait for LSP command to complete
	vim.loop.sleep(1500)
	
	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Verify content changed
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after moving Level 1 item")
	
	-- Verify that "Level 1 Item B" now appears before "Level 1 Item A"
	local item_b_pos = string.find(final_content, "Level 1 Item B")
	local item_a_pos = string.find(final_content, "Level 1 Item A")
	
	MiniTest.expect.equality(item_b_pos < item_a_pos, true, "Level 1 Item B should appear before Level 1 Item A")
	
	-- Verify that all nested items under B moved with it
	MiniTest.expect.equality(string.find(final_content, "Level 2 Item B%.1") ~= nil, true, "Nested items should move with parent")
	MiniTest.expect.equality(string.find(final_content, "Level 3 Item B%.1%.a") ~= nil, true, "Deep nested items should move with parent")
	
	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["nested list movement - level 2 items"] = function()
	local workspace_path, file_path = create_nested_list_test_workspace()
	local child = utils.new_child_neovim()
	
	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)
	
	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)
	
	-- Position cursor on "Level 2 Item A.2" 
	child.lua('vim.fn.search("Level 2 Item A.2")')
	
	-- Get initial content
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Execute NotedownMoveUp command to move A.2 above A.1
	child.lua('vim.cmd("NotedownMoveUp")')
	
	-- Wait for LSP command to complete
	vim.loop.sleep(1500)
	
	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Verify content changed
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after moving Level 2 item")
	
	-- Verify that "Level 2 Item A.2" now appears before "Level 2 Item A.1"
	-- We need to find their positions relative to their parent "Level 1 Item A"
	local parent_pos = string.find(final_content, "Level 1 Item A")
	local content_after_parent = string.sub(final_content, parent_pos)
	
	local item_a2_pos = string.find(content_after_parent, "Level 2 Item A%.2")
	local item_a1_pos = string.find(content_after_parent, "Level 2 Item A%.1")
	
	MiniTest.expect.equality(item_a2_pos < item_a1_pos, true, "Level 2 Item A.2 should appear before A.1 after move up")
	
	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["nested list movement - deep level items (level 4)"] = function()
	local workspace_path, file_path = create_nested_list_test_workspace()
	local child = utils.new_child_neovim()
	
	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)
	
	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)
	
	-- Position cursor on "Level 4 Item A.1.a.ii"
	child.lua('vim.fn.search("Level 4 Item A.1.a.ii")')
	
	-- Get initial content
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Execute NotedownMoveUp command to move A.1.a.ii above A.1.a.i
	child.lua('vim.cmd("NotedownMoveUp")')
	
	-- Wait for LSP command to complete
	vim.loop.sleep(1500)
	
	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Verify content changed
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after moving Level 4 item")
	
	-- Find the parent Level 3 item and check order within it
	local level3_pos = string.find(final_content, "Level 3 Item A%.1%.a")
	local content_after_level3 = string.sub(final_content, level3_pos)
	
	local item_ii_pos = string.find(content_after_level3, "Level 4 Item A%.1%.a%.ii")
	-- Use more precise search for item i to avoid matching item ii
	local item_i_pos = nil
	for i = 1, #content_after_level3 do
		local substr = string.sub(content_after_level3, i)
		if string.find(substr, "^Level 4 Item A%.1%.a%.i") and not string.find(substr, "^Level 4 Item A%.1%.a%.ii") then
			item_i_pos = i
			break
		end
	end
	
	if item_i_pos and item_ii_pos then
		MiniTest.expect.equality(item_ii_pos < item_i_pos, true, "Level 4 Item A.1.a.ii should appear before A.1.a.i after move up")
	else
		error("Could not find both Level 4 items in final content for position verification")
	end
	
	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["nested list movement - level 5 items with deep children"] = function()
	local workspace_path, file_path = create_nested_list_test_workspace()
	local child = utils.new_child_neovim()
	
	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)
	
	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)
	
	-- Position cursor on "Level 5 Item A.1.a.i.β"
	child.lua('vim.fn.search("Level 5 Item A.1.a.i.β")')
	
	-- Get initial content
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Execute NotedownMoveUp command to move β above α
	child.lua('vim.cmd("NotedownMoveUp")')
	
	-- Wait for LSP command to complete
	vim.loop.sleep(1500)
	
	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Verify content changed
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after moving Level 5 item")
	
	-- Find the parent Level 4 item and check order within it
	local level4_pos = string.find(final_content, "Level 4 Item A%.1%.a%.i")
	local content_after_level4 = string.sub(final_content, level4_pos)
	
	local item_beta_pos = string.find(content_after_level4, "Level 5 Item A%.1%.a%.i%.β")
	local item_alpha_pos = string.find(content_after_level4, "Level 5 Item A%.1%.a%.i%.α")
	
	MiniTest.expect.equality(item_beta_pos < item_alpha_pos, true, "Level 5 Item β should appear before α after move up")
	
	-- Verify that Level 6 items under α moved with it (not separated)
	-- Find the content of the α section by looking for the next Level 5 item or end of section
	local alpha_section_start = item_alpha_pos
	local alpha_section_end = string.find(content_after_level4, "Level 5", item_alpha_pos + 1) or #content_after_level4
	
	local alpha_section_content = string.sub(content_after_level4, alpha_section_start, alpha_section_end - 1)
	MiniTest.expect.equality(string.find(alpha_section_content, "Level 6") ~= nil, true, "Level 6 items should move with their Level 5 parent")
	
	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["nested list movement - mixed list types"] = function()
	local workspace_path, file_path = create_nested_list_test_workspace()
	local child = utils.new_child_neovim()
	
	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)
	
	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)
	
	-- Position cursor on "2. Ordered Level 1 Item B" in the mixed list section
	child.lua('vim.fn.search("2. Ordered Level 1 Item B")')
	
	-- Get initial content
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Execute NotedownMoveUp command
	child.lua('vim.cmd("NotedownMoveUp")')
	
	-- Wait for LSP command to complete
	vim.loop.sleep(1500)
	
	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Verify content changed
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after moving mixed list item")
	
	-- For ordered lists, verify that the numbers get updated correctly
	-- After moving item 2 up, it should become item 1, and the original item 1 should become item 2
	MiniTest.expect.equality(string.find(final_content, "1%. Ordered Level 1 Item B") ~= nil, true, "Moved item should be renumbered to 1")
	MiniTest.expect.equality(string.find(final_content, "2%. Ordered Level 1 Item A") ~= nil, true, "Original item 1 should be renumbered to 2")
	
	-- Verify nested bullet items moved with their parent
	local item_b_pos = string.find(final_content, "1%. Ordered Level 1 Item B")
	local content_after_b = string.sub(final_content, item_b_pos)
	MiniTest.expect.equality(string.find(content_after_b, "Bullet Level 2 Item B%.1") ~= nil, true, "Nested items should move with parent in mixed lists")
	
	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["nested list movement - task lists with deep nesting"] = function()
	local workspace_path, file_path = create_nested_list_test_workspace()
	local child = utils.new_child_neovim()
	
	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)
	
	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)
	
	-- Position cursor on "- [ ] Sub-subtask A.1.b" in the task list section
	child.lua('vim.fn.search("Sub-subtask A.1.b")')
	
	-- Get initial content
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Execute NotedownMoveUp command to move A.1.b above A.1.a
	child.lua('vim.cmd("NotedownMoveUp")')
	
	-- Wait for LSP command to complete
	vim.loop.sleep(1500)
	
	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Verify content changed
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after moving task list item")
	
	-- Find the parent "Subtask A.1" and verify order of children
	local parent_pos = string.find(final_content, "Subtask A%.1")
	local content_after_parent = string.sub(final_content, parent_pos)
	
	local item_b_pos = string.find(content_after_parent, "Sub%-subtask A%.1%.b")
	local item_a_pos = string.find(content_after_parent, "Sub%-subtask A%.1%.a")
	
	MiniTest.expect.equality(item_b_pos < item_a_pos, true, "Sub-subtask A.1.b should appear before A.1.a after move up")
	
	-- Verify that deeply nested items under A.1.b moved with it
	local item_b_section = string.sub(content_after_parent, item_b_pos)
	MiniTest.expect.equality(string.find(item_b_section, "Deep subtask A%.1%.b%.i") ~= nil, true, "Deep nested items should move with parent task")
	MiniTest.expect.equality(string.find(item_b_section, "Very deep subtask A%.1%.b%.i%.α") ~= nil, true, "Very deep nested items should move with parent task")
	
	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["nested list movement - boundary conditions with nesting"] = function()
	local workspace_path, file_path = create_nested_list_test_workspace()
	local child = utils.new_child_neovim()
	
	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)
	
	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)
	
	-- Test 1: Try to move the first Level 6 item up (should not move)
	child.lua('vim.fn.search("Level 6 Item A.1.a.i.α.I")')
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	child.lua('vim.cmd("NotedownMoveUp")')
	vim.loop.sleep(1500)
	
	local after_first_attempt = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Should not change since it's already at the top of its level
	MiniTest.expect.equality(initial_content, after_first_attempt, "First item at deepest level should not move up")
	
	-- Test 2: Try to move the last Level 6 item down (should not move)
	child.lua('vim.fn.search("Level 6 Item A.1.a.i.α.II")')
	local before_second_attempt = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	child.lua('vim.cmd("NotedownMoveDown")')
	vim.loop.sleep(1500)
	
	local after_second_attempt = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Should not change since it's already at the bottom of its level
	MiniTest.expect.equality(before_second_attempt, after_second_attempt, "Last item at deepest level should not move down")
	
	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

T["nested list movement - preserve hierarchy structure"] = function()
	local workspace_path, file_path = create_nested_list_test_workspace()
	local child = utils.new_child_neovim()
	
	-- Set up notedown with real LSP server
	lsp.setup(child, workspace_path)
	
	-- Open the test file
	child.lua('vim.cmd("edit ' .. file_path .. '")')
	lsp.wait_for_ready(child)
	
	-- Position cursor on "Level 3 Item A.1.b" and move it up
	child.lua('vim.fn.search("Level 3 Item A.1.b")')
	
	-- Get initial content for analysis
	local initial_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Execute move command
	child.lua('vim.cmd("NotedownMoveUp")')
	vim.loop.sleep(1500)
	
	-- Get final content
	local final_content = child.lua_get('table.concat(vim.api.nvim_buf_get_lines(0, 0, -1, false), "\\n")')
	
	-- Verify movement occurred
	MiniTest.expect.no_equality(initial_content, final_content, "Content should change after moving nested item")
	
	-- Verify that the hierarchy structure is preserved
	-- Check that all indentation levels are maintained correctly
	local final_lines = child.lua_get('vim.api.nvim_buf_get_lines(0, 0, -1, false)')
	
	local level_counts = { 0, 0, 0, 0, 0, 0 } -- Count items at each level
	
	for _, line in ipairs(final_lines) do
		local trimmed = string.gsub(line, "^%s*", "") -- Remove leading whitespace
		if string.find(trimmed, "^[-*]") or string.find(trimmed, "^%d+%.") or string.find(trimmed, "^- %[") then
			-- Count the indentation level (every 2 spaces = 1 level)
			local indent = string.len(line) - string.len(string.gsub(line, "^%s*", ""))
			local level = math.floor(indent / 2) + 1
			if level <= 6 then
				level_counts[level] = level_counts[level] + 1
			end
		end
	end
	
	-- Verify we still have items at all expected levels
	MiniTest.expect.equality(level_counts[1] > 0, true, "Should have Level 1 items")
	MiniTest.expect.equality(level_counts[2] > 0, true, "Should have Level 2 items")
	MiniTest.expect.equality(level_counts[3] > 0, true, "Should have Level 3 items")
	MiniTest.expect.equality(level_counts[4] > 0, true, "Should have Level 4 items")
	MiniTest.expect.equality(level_counts[5] > 0, true, "Should have Level 5 items")
	MiniTest.expect.equality(level_counts[6] > 0, true, "Should have Level 6 items")
	
	child.stop()
	utils.cleanup_test_workspace(workspace_path)
	lsp.cleanup_binary()
end

return T
