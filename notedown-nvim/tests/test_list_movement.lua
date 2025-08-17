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

-- Golden file tests for list item movement functionality

local MiniTest = require("mini.test")
local golden = require("helpers.golden")
local lsp_shared = require("helpers.lsp_shared")

local T = MiniTest.new_set()

-- Initialize shared LSP session once for the entire test suite
lsp_shared.initialize()

-- Register cleanup function to run when all tests complete
_G._notedown_list_movement_cleanup = lsp_shared.cleanup

-- ========================================
-- SIMPLE LIST MOVEMENT TESTS
-- ========================================

T["simple - move second item down"] = function()
	golden.test_list_movement("simple", "move_second_down", {
		search_pattern = "Second item",
		command = "NotedownMoveDown",
		expected_cursor = { 5, 2 }, -- "Second item" moves from line 4 to line 5, cursor at "S"
	})
end

T["simple - move third item up"] = function()
	golden.test_list_movement("simple", "move_third_up", {
		search_pattern = "Third item",
		command = "NotedownMoveUp",
		expected_cursor = { 4, 2 }, -- "Third item" moves from line 5 to line 4, cursor at "T"
	})
end

T["simple - move first item up (no change)"] = function()
	golden.test_list_movement("simple", "move_first_up_no_change", {
		search_pattern = "First item",
		command = "NotedownMoveUp",
		expected_cursor = { 3, 2 }, -- No movement, cursor stays at original position
	})
end

T["simple - move fourth item down (no change)"] = function()
	golden.test_list_movement("simple", "move_fourth_down_no_change", {
		search_pattern = "Fourth item",
		command = "NotedownMoveDown",
		expected_cursor = { 6, 2 }, -- No movement, cursor stays at original position
	})
end

-- ========================================
-- NESTED LIST MOVEMENT TESTS
-- ========================================

T["nested - move Level 1 Item B up"] = function()
	golden.test_list_movement("nested", "level1_move_b_up", {
		search_pattern = "Level 1 Item B",
		command = "NotedownMoveUp",
		expected_cursor = { 3, 2 }, -- "Level 1 Item B" moves to line 3, cursor at "L"
	})
end

T["nested - move Level 2 Item A.2 up"] = function()
	golden.test_list_movement("nested", "level2_move_a2_up", {
		search_pattern = "Level 2 Item A.2",
		command = "NotedownMoveUp",
		expected_cursor = { 4, 4 }, -- "Level 2 Item A.2" moves to line 4, cursor at "L"
	})
end

T["nested - move Level 4 Item A.1.a.ii up"] = function()
	golden.test_list_movement("nested", "level4_move_ii_up", {
		search_pattern = "Level 4 Item A.1.a.ii",
		command = "NotedownMoveUp",
		expected_cursor = { 6, 8 }, -- "Level 4 Item A.1.a.ii" moves to line 6, cursor at "L"
	})
end

T["nested - move Level 5 Item A.1.a.i.β up"] = function()
	golden.test_list_movement("nested", "level5_move_beta_up", {
		search_pattern = "Level 5 Item A.1.a.i.β",
		command = "NotedownMoveUp",
		expected_cursor = { 7, 10 }, -- "Level 5 Item A.1.a.i.β" moves to line 7, cursor at "L"
	})
end

T["nested - mixed list renumbering"] = function()
	golden.test_list_movement("nested", "mixed_list_renumber", {
		search_pattern = "2. Ordered Level 1 Item B",
		command = "NotedownMoveUp",
		expected_cursor = { 27, 3 }, -- "Ordered Level 1 Item B" moves to line 27, cursor at "O"
	})
end

-- ========================================
-- TASK LIST MOVEMENT TESTS
-- ========================================

T["tasks - move Sub-subtask A.1.b up"] = function()
	golden.test_list_movement("tasks", "move_subtask_a1b_up", {
		search_pattern = "Sub-subtask A.1.b",
		command = "NotedownMoveUp",
		expected_cursor = { 5, 10 }, -- Actual cursor position
	})
end

-- ========================================
-- BOUNDARY CONDITION TESTS
-- ========================================

T["nested - try to move first Level 6 item up (no change)"] = function()
	golden.test_list_movement("nested", "boundary_level6_first_up_no_change", {
		search_pattern = "Level 6 Item A.1.a.i.α.I",
		command = "NotedownMoveUp",
		expected_cursor = { 8, 12 }, -- No movement, cursor stays at original position
	})
end

return T
