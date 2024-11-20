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

type Fetcher = collections.Fetcher[TaskClient, Task]
type ListOption = collections.ListOption[Task]

func (c *TaskClient) TaskSummary() int {
	tasks := 0
	c.tasksMutex.RLock()
	defer c.tasksMutex.RUnlock()
	for _, doc := range c.tasks {
		tasks += len(doc)
	}
	return tasks
}

// Opts are applied in order so filters should be applied before sorters
func (c *TaskClient) ListTasks(fetcher Fetcher, opts ...ListOption) []Task {
	return collections.List(c, fetcher, opts...)
}
