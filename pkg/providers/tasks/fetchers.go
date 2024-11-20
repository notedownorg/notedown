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

import "github.com/notedownorg/notedown/pkg/providers/pkg/collections"

func FetchAllTasks() collections.Fetcher[TaskClient, Task] {
	return func(c *TaskClient) []Task {
		var tasks []Task
		c.tasksMutex.RLock()
		for _, document := range c.tasks {
			for _, task := range document {
				tasks = append(tasks, task)
			}
		}
		c.tasksMutex.RUnlock()
		return tasks
	}
}

func FetchTasksForDocument(document string) collections.Fetcher[TaskClient, Task] {
	return func(c *TaskClient) []Task {
		var tasks []Task
		c.tasksMutex.RLock()
		for _, task := range c.tasks[document] {
			tasks = append(tasks, task)
		}
		c.tasksMutex.RUnlock()
		return tasks
	}
}
