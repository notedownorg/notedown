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
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/documents/reader"
)

type TaskFetcher func(c *Client) []ast.Task

func FetchAllTasks() TaskFetcher {
	return func(c *Client) []ast.Task {
		var tasks []ast.Task
		c.tasksMutex.RLock()
		for _, document := range c.tasks {
			for _, task := range document {
				tasks = append(tasks, *task)
			}
		}
		c.tasksMutex.RUnlock()
		return tasks
	}
}

func FetchTasksForDocument(document string) TaskFetcher {
	return func(c *Client) []ast.Task {
		var tasks []ast.Task
		c.tasksMutex.RLock()
		for _, task := range c.tasks[document] {
			tasks = append(tasks, *task)
		}
		c.tasksMutex.RUnlock()
		return tasks
	}
}

type DocumentFetcher func(c *Client) map[string]reader.Document

func FetchAllDocuments() DocumentFetcher {
	return func(c *Client) map[string]reader.Document {
		c.documentsMutex.RLock()
		defer c.documentsMutex.RUnlock()
		documents := make(map[string]reader.Document)
		for path, document := range c.documents {
			documents[path] = *document
		}
		return documents
	}
}
