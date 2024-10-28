// Copyright 2024 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tasks

import (
	"time"

	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/documents/reader"
)

type TaskFilter func(ast.Task) bool

func WithFilters(filters ...TaskFilter) ListTasksOptions {
	return func(tasks []ast.Task) []ast.Task {
		for _, filter := range filters {
			tasks = filterTasks(tasks, filter)
		}
		return tasks
	}
}

// Priorities are OR'd together because a task can't have multiple priorities.
func FilterByPriority(priority ...int) TaskFilter {
	return func(task ast.Task) bool {
		for _, p := range priority {
			taskPriority := task.Priority()
			if taskPriority != nil && *taskPriority == p {
				return true
			}
		}
		return false
	}
}

// Statuses are OR'd together because a task can only have one status.
func FilterByStatus(status ...ast.Status) TaskFilter {
	return func(task ast.Task) bool {
		for _, s := range status {
			if task.Status() == s {
				return true
			}
		}
		return false
	}
}

// Following Go's time package, after and before are inclusive (include equal to).
func FilterByDueDate(after *time.Time, before *time.Time) TaskFilter {
	return func(t ast.Task) bool {
		if t.Due() == nil {
			return false
		}
		if after != nil && t.Due().Before(*after) {
			return false
		}
		if before != nil && t.Due().After(*before) {
			return false
		}
		return true
	}
}

func FilterByCompletedDate(after *time.Time, before *time.Time) TaskFilter {
	return func(t ast.Task) bool {
		if t.Completed() == nil {
			return false
		}
		if after != nil && t.Completed().Before(*after) {
			return false
		}
		if before != nil && t.Completed().After(*before) {
			return false
		}
		return true
	}
}

type DocumentFilter func(path string, document reader.Document) bool

func FilterByDocumentType(documentType string) DocumentFilter {
	return func(_ string, document reader.Document) bool {
		return document.Metadata["type"] == documentType
	}
}
