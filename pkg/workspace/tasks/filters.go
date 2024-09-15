package tasks

import "github.com/liamawhite/nl/pkg/ast"

type TaskFilter func(Task) bool

// Priorities are OR'd together because a task can only have one priority
func FilterByPriority(priority ...int) TaskFilter {
	return func(task Task) bool {
		for _, p := range priority {
			if task.Priority != nil && *task.Priority == p {
				return true
			}
		}
		return false
	}
}

// Statuses are OR'd together because a task can only have one status
func FilterByStatus(status ...ast.Status) TaskFilter {
	return func(task Task) bool {
		for _, s := range status {
			if task.Status == s {
				return true
			}
		}
		return false
	}
}
