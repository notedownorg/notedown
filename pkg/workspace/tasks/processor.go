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

func (c *Client) processDocuments(feed <-chan reader.Event) {
	for {
		select {
		case event := <-feed:
			switch event.Op {
			case reader.Delete:
				c.documentsMutex.Lock()
				delete(c.documents, event.Key)
				c.documentsMutex.Unlock()
				c.tasksMutex.Lock()
				delete(c.tasks, event.Key)
				c.tasksMutex.Unlock()
				c.events <- Event{Op: Delete}
			case reader.Change:
				c.handleChanges(event)
				c.events <- Event{Op: Change}
			case reader.Load:
				c.handleChanges(event)
				c.events <- Event{Op: Load}
			case reader.SubscriberLoadComplete:
				c.initialLoadComplete = true
			}

		}
	}
}

func (c *Client) handleChanges(event reader.Event) {
	tasks := make(map[int]*ast.Task)
	for i := range event.Document.Tasks {
		task := event.Document.Tasks[i]
		tasks[task.Line()] = &task
	}
	c.tasksMutex.Lock()
	c.tasks[event.Key] = tasks
	c.tasksMutex.Unlock()
	c.documentsMutex.Lock()
	c.documents[event.Key] = &event.Document
	c.documentsMutex.Unlock()
}
