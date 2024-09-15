package tasks_test

import (
	"time"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/workspace/documents"
	"github.com/liamawhite/nl/pkg/workspace/tasks"
)

func buildClient(events ...documents.Event) (*tasks.Client, chan documents.Event) {
	feed := make(chan documents.Event)
	client := tasks.NewClient(feed)
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

func intPtr(i int) *int {
	return &i
}

func tasksBuilder(doc documents.Document) []tasks.Task {
	res := make([]tasks.Task, len(doc.Tasks))
	for i, t := range doc.Tasks {
		res[i] = toTask(t, doc.Hash)
	}
	return res
}

func toTask(t ast.Task, documentHash string) tasks.Task {
	return tasks.Task{Task: t, DocumentHash: documentHash}
}

func defaultEvents() []documents.Event {
	return []documents.Event{
		{
			Op:  documents.Change,
			Key: "one.md",
			Document: documents.Document{
				Document: ast.Document{
					Tasks: []ast.Task{
						{
							Name:   "Task 1",
							Line:   1,
							Status: ast.Todo,
						},
						{
							Name:     "Task 2",
							Line:     2,
							Priority: intPtr(1),
							Status:   ast.Doing,
						},
					},
				},
			},
		},
		{
			Op:  documents.Change,
			Key: "two.md",
			Document: documents.Document{
				Document: ast.Document{
					Tasks: []ast.Task{
						{
							Name:     "Task 3",
							Line:     1,
							Priority: intPtr(2),
							Status:   ast.Done,
						},
						{
							Name:     "Task 4",
							Line:     2,
							Priority: intPtr(3),
							Status:   ast.Abandoned,
						},
						{
							Name:     "Task 5",
							Line:     3,
							Priority: intPtr(10),
							Status:   ast.Blocked,
						},
					},
				},
			},
		},
	}
}
