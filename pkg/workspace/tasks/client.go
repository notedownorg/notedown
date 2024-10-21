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
	"sync"

	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/documents/reader"
	"github.com/notedownorg/notedown/pkg/workspace/documents/writer"
)

type Client struct {
	// cache maps between file paths and line numbers to tasks it should ONLY be updated in response
	// to events from the docuuments client and should otherwise be read-only.
	cache map[string]map[int]*ast.Task
	mutex sync.RWMutex

	initialLoaderWaiter *sync.WaitGroup

	writer writer.LineWriter
}

type clientOptions func(*Client)

func WithInitialLoadWaiter(wg *sync.WaitGroup) clientOptions {
	return func(client *Client) {
		client.initialLoaderWaiter = wg
	}
}

func NewClient(writer writer.LineWriter, feed <-chan reader.Event, opts ...clientOptions) *Client {
	client := &Client{
		cache:  make(map[string]map[int]*ast.Task),
		writer: writer,
	}
	go client.processDocuments(feed)

	for _, opt := range opts {
		opt(client)
	}

	if client.initialLoaderWaiter != nil {
		client.initialLoaderWaiter.Wait()
	}

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
			case reader.Change, reader.Load:
				if event.Document.Tasks == nil || len(event.Document.Tasks) == 0 {
					break
				}
				tasks := make(map[int]*ast.Task)
				for i := range event.Document.Tasks {
					task := event.Document.Tasks[i]
					tasks[task.Line()] = &task
				}
				c.mutex.Lock()
				c.cache[event.Key] = tasks
				c.mutex.Unlock()

				if event.Op == reader.Load && c.initialLoaderWaiter != nil {
					c.initialLoaderWaiter.Done()
				}
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
