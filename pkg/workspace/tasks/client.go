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
	"time"

	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/documents/reader"
	"github.com/notedownorg/notedown/pkg/workspace/documents/writer"
)

type Client struct {
	// tasks maps between file paths and line numbers to tasks it should ONLY be updated in response
	// to events from the docuuments client and should otherwise be read-only.
	tasks      map[string]map[int]*ast.Task
	tasksMutex sync.RWMutex

	// documents maps between file paths and documents it should ONLY be updated in response
	// to events from the docuuments client and should otherwise be read-only.
	documents      map[string]*reader.Document
	documentsMutex sync.RWMutex

	subscribers []chan Event
	events      chan Event

	initialLoadComplete bool

	writer writer.LineWriter
}

type clientOptions func(*Client)

// Inform NewClient to wait for the initial load to complete before returning
func WithInitialLoadWaiter(tick time.Duration) clientOptions {
	return func(client *Client) {
		for !client.initialLoadComplete {
			time.Sleep(tick)
		}
	}
}

func NewClient(writer writer.LineWriter, feed <-chan reader.Event, opts ...clientOptions) *Client {
	client := &Client{
		tasks:               make(map[string]map[int]*ast.Task),
		documents:           make(map[string]*reader.Document),
		writer:              writer,
		events:              make(chan Event),
		subscribers:         make([]chan Event, 0),
		initialLoadComplete: false,
	}

	go client.processDocuments(feed)
	go client.eventDispatcher()

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Summary() (tasks int, documents int) {
	c.tasksMutex.RLock()
	defer c.tasksMutex.RUnlock()
	for _, doc := range c.tasks {
		tasks += len(doc)
		documents++
	}
	return tasks, documents
}

func (c *Client) ListDocuments(fetcher DocumentFetcher, filters ...DocumentFilter) map[string]reader.Document {
	documents := fetcher(c)
	for _, filter := range filters {
		documents = filterDocuments(documents, filter)
	}
	return documents
}

func filterDocuments(documents map[string]reader.Document, filter DocumentFilter) map[string]reader.Document {
	filtered := make(map[string]reader.Document)
	for path, document := range documents {
		if filter(path, document) {
			filtered[path] = document
		}
	}
	return filtered
}

type ListTasksOptions func(tasks []ast.Task) []ast.Task

// Opts are applied in order so filters should be applied before sorters
func (c *Client) ListTasks(fetcher TaskFetcher, opts ...ListTasksOptions) []ast.Task {
	tasks := fetcher(c)

	for _, opt := range opts {
		tasks = opt(tasks)
	}

	return tasks
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
