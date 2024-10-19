package tasks

import (
	"fmt"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/workspace/documents/writer"
)

func (c *Client) Create(path string, name string, status ast.Status, options ...ast.TaskOption) error {
	// TODO:
	// This won't work with adding subtasks or any other tasks on specific lines as the writer will reject them with an empty hash
	// Probably need to have some kind of with parent task option for subtasks?
	// Just reject any tasks with a line number?
	task := ast.NewTask(ast.NewIdentifier(path, ""), name, status, options...)
	err := c.writer.AddLine(writer.Document{Path: path}, task.Line(), task)
	if err != nil {
		return fmt.Errorf("failed to add task: %v: %w", task, err)
	}
	return nil
}

func (c *Client) Update(t ast.Task) error {
	err := c.writer.UpdateLine(writer.Document{Path: t.Path(), Hash: t.Version()}, t.Line(), t)
	if err != nil {
		return fmt.Errorf("failed to update task: %v: %w", t, err)
	}
	return nil
}

func (c *Client) Delete(t ast.Task) error {
	err := c.writer.RemoveLine(writer.Document{Path: t.Path(), Hash: t.Version()}, t.Line())
	if err != nil {
		return fmt.Errorf("failed to remove task: %v: %w", t, err)
	}
	return nil
}

// Move moves a task to the end of the file at the specified path
func (c *Client) Move(t ast.Task, path string) error {
	// Do the add first so we don't accidentally remove the task without adding it to the new file
	err := c.writer.AddLine(writer.Document{Path: path}, writer.AtEnd, t)
	if err != nil {
		return fmt.Errorf("failed to add new task when moving: %v: %w", t, err)
	}
	err = c.writer.RemoveLine(writer.Document{Path: t.Path(), Hash: t.Version()}, t.Line())
	if err != nil {
		return fmt.Errorf("failed to remove existing task when moving: %v: %w", t, err)
	}
	return nil
}
