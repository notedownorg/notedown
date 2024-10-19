package tasks_test

import (
	"time"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/workspace/documents/reader"
	"github.com/liamawhite/nl/pkg/workspace/tasks"
)

func buildClient(events []reader.Event, validators ...validator) (*tasks.Client, chan reader.Event) {
	feed := make(chan reader.Event)
	client := tasks.NewClient(&MockLineWriter{validators: validators}, feed)
	for _, event := range events {
		feed <- event
	}

	// Wait for the events to be processed
	// NOTE: This assumes that all events are document creation events
	for len(client.ListDocuments()) != len(events) {
		time.Sleep(100 * time.Millisecond)
	}

	return client, feed
}

func defaultEvents() []reader.Event {
	return []reader.Event{
		{
			Op:  reader.Change,
			Key: "one.md",
			Document: reader.Document{
				Document: ast.Document{
					Tasks: []ast.Task{
						ast.NewTask(ast.NewIdentifier("one.md", "version"), "Task 1", ast.Todo, ast.WithLine(1)),
						ast.NewTask(ast.NewIdentifier("one.md", "version"), "Task 2", ast.Doing, ast.WithLine(2), ast.WithPriority(1)),
					},
				},
			},
		},
		{
			Op:  reader.Change,
			Key: "two.md",
			Document: reader.Document{
				Document: ast.Document{
					Tasks: []ast.Task{
						ast.NewTask(ast.NewIdentifier("two.md", "version"), "Task 3", ast.Done, ast.WithLine(1), ast.WithPriority(2)),
						ast.NewTask(ast.NewIdentifier("two.md", "version"), "Task 4", ast.Abandoned, ast.WithLine(2), ast.WithPriority(3)),
						ast.NewTask(ast.NewIdentifier("two.md", "version"), "Task 5", ast.Blocked, ast.WithLine(3), ast.WithPriority(10)),
					},
				},
			},
		},
	}
}
