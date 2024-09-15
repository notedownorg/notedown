package tasks

import "github.com/liamawhite/nl/pkg/ast"

type TaskFilter func(ast.Task) bool

// Priorities are OR'd together
func FilterByPriority(priority ...int) TaskFilter {
	return func(task ast.Task) bool {
		for _, p := range priority {
			if task.Priority != nil && *task.Priority == p {
				return true
			}
		}
		return false
	}
}

func FilterByStatus(status ...ast.Status) TaskFilter {
	return func(task ast.Task) bool {
		for _, s := range status {
			if task.Status == s {
				return true
			}
		}
		return false
	}
}
