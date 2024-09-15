package tasks

import (
	"sync"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/workspace/documents/reader"
)

type Client struct {
	// cache maps between file paths and line numbers to tasks it should ONLY be updated in response
	// to events from the docuuments client and should otherwise be read-only.
	cache map[string]map[int]*ast.Task
	mutex sync.RWMutex
}

func NewClient(feed <-chan reader.Event) *Client {
	client := &Client{
		cache: make(map[string]map[int]*ast.Task),
	}
	go client.processDocuments(feed)
	return client
}

func (c *Client) processDocuments(feed <-chan reader.Event) {
	for {
		select {
		case event := <-feed:
			switch event.Op {
			case reader.Delete:
				c.mutex.Lock()
				delete(c.cache, event.Key)
				c.mutex.Unlock()
			case reader.Change:
				if event.Document.Tasks == nil || len(event.Document.Tasks) == 0 {
					break
				}
				tasks := make(map[int]*ast.Task)
				for i := range event.Document.Tasks {
					task := event.Document.Tasks[i]
					tasks[task.Line] = &task
				}
				c.mutex.Lock()
				c.cache[event.Key] = tasks
				c.mutex.Unlock()
			}

		}
	}
}

func (c *Client) ListDocuments() []string {
	var documents []string
	c.mutex.RLock()
	for document := range c.cache {
		documents = append(documents, document)
	}
	c.mutex.RUnlock()
	return documents
}

func (c *Client) ListTasks(fetcher TaskFetcher, filters ...TaskFilter) ([]ast.Task, error) {
	tasks, err := fetcher(c)
	if err != nil {
		return nil, err
	}

	for _, filter := range filters {
		tasks = filterTasks(tasks, filter)
	}

	return tasks, nil
}
func filterTasks(tasks []ast.Task, filter TaskFilter) []ast.Task {
	var filtered []ast.Task
	for _, task := range tasks {
		if filter(task) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}
