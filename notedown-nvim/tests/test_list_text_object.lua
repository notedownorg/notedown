-- Copyright 2025 Notedown Authors
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

-- Tests for list text object functionality

local MiniTest = require("mini.test")
local golden = require("helpers.golden")
local lsp_shared = require("helpers.lsp_shared")

local T = MiniTest.new_set()

-- Initialize shared LSP session once for the entire test suite
lsp_shared.initialize()

-- Register cleanup function to run when all tests complete
_G._notedown_list_text_object_cleanup = lsp_shared.cleanup

-- ========================================
-- SIMPLE LIST TEXT OBJECT TESTS
-- ========================================

T["simple - yank first item"] = function()
	golden.test_text_object("simple", "yank_first_item", {
		search_pattern = "First item",
		operation = "yal", -- yank around list
		expected_register_content = "- First item\n",
	})
end

T["simple - yank second item"] = function()
	golden.test_text_object("simple", "yank_second_item", {
		search_pattern = "Second item",
		operation = "yal",
		expected_register_content = "- Second item\n",
	})
end

T["simple - delete third item"] = function()
	golden.test_text_object("simple", "delete_third_item", {
		search_pattern = "Third item",
		operation = "dal", -- delete around list
		expected_cursor = { 5, 2 }, -- Cursor should be on "Fourth item"
	})
end

-- ========================================
-- NESTED LIST TEXT OBJECT TESTS
-- ========================================

T["nested - yank level 1 item with children"] = function()
	golden.test_text_object("nested", "yank_level1_with_children", {
		search_pattern = "Level 1 Item A",
		operation = "yal",
		expected_register_content = "- Level 1 Item A\n  - Level 2 Item A.1\n    - Level 3 Item A.1.a\n      - Level 4 Item A.1.a.i\n        - Level 5 Item A.1.a.i.α\n          - Level 6 Item A.1.a.i.α.I\n      - Level 4 Item A.1.a.ii\n        - Level 5 Item A.1.a.i.β\n  - Level 2 Item A.2\n",
	})
end

T["nested - delete level 4 item with children"] = function()
	golden.test_text_object("nested", "delete_level4_with_children", {
		search_pattern = "Level 4 Item A.1.a.i",
		operation = "dal",
		expected_cursor = { 7, 8 }, -- Should move to next item
	})
end

T["nested - yank deep nested single item"] = function()
	golden.test_text_object("nested", "yank_level6_single", {
		search_pattern = "Level 6 Item A.1.a.i.α.I",
		operation = "yal",
		expected_register_content = "          - Level 6 Item A.1.a.i.α.I\n",
	})
end

-- ========================================
-- TASK LIST TEXT OBJECT TESTS
-- ========================================

T["tasks - yank completed task with subtasks"] = function()
	golden.test_text_object("tasks", "yank_completed_with_subtasks", {
		search_pattern = "- [x] Completed task",
		operation = "yal",
		expected_register_content = "- [x] Completed task\n  - [ ] Subtask A\n    - [ ] Sub-subtask A.1.a\n    - [x] Sub-subtask A.1.b\n  - [x] Subtask B\n",
	})
end

T["tasks - delete incomplete subtask"] = function()
	golden.test_text_object("tasks", "delete_incomplete_subtask", {
		search_pattern = "Subtask A",
		operation = "dal",
		expected_cursor = { 6, 4 }, -- Move to next subtask
	})
end

-- ========================================
-- BOUNDARY CONDITION TESTS
-- ========================================

T["simple - no list item at cursor"] = function()
	golden.test_text_object("simple", "no_list_item", {
		search_pattern = "Some text after", -- Not a list item
		operation = "yal",
		should_fail = true,
		expected_warning = "No list item found at cursor",
	})
end

T["nested - mixed list types"] = function()
	golden.test_text_object("nested", "mixed_list_numbered", {
		search_pattern = "1. Ordered Level 1 Item A",
		operation = "yal",
		expected_register_content = "1. Ordered Level 1 Item A\n    a. Ordered Level 2 Item A.a\n    b. Ordered Level 2 Item A.b\n",
	})
end

-- ========================================
-- MULTI-LINE CONTENT TESTS
-- ========================================

T["multiline - list item with content"] = function()
	golden.test_text_object("multiline", "item_with_content", {
		search_pattern = "List item with",
		operation = "yal",
		expected_register_content = "- List item with\n  additional content\n  spanning multiple lines\n",
	})
end

return T
