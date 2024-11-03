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
	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
)

func WithSorters(sorters ...collections.Sorter[Task]) collections.ListOption[Task] {
	return collections.Sort(collections.FallthroughDeterministic(sorters...))
}

func SortByPriority() collections.Sorter[Task] {
	return func(a, b Task) int {
		if a.Priority() == nil && b.Priority() == nil {
			return 0
		}
		if a.Priority() == nil {
			return 1
		}
		if b.Priority() == nil {
			return -1
		}
		if *a.Priority() == *b.Priority() {
			return 0
		}
		if *a.Priority() < *b.Priority() {
			return -1
		}
		return 1
	}
}

func AgendaOrder() (Status, Status, Status, Status, Status) {
	return Doing, Todo, Blocked, Done, Abandoned
}

func KanbanOrder() (Status, Status, Status, Status, Status) {
	return Todo, Blocked, Doing, Done, Abandoned
}

func SortByStatus(first, second, third, fourth, fifth Status) collections.Sorter[Task] {
	return func(a, b Task) int {
		if a.Status() == b.Status() {
			return 0
		}
		switch a.Status() {
		case first:
			if b.Status() == first {
				return 0
			}
			return -1
		case second:
			if b.Status() == first {
				return 1
			}
			return -1
		case third:
			if b.Status() == first || b.Status() == second {
				return 1
			}
			return -1
		case fourth:
			if b.Status() == first || b.Status() == second || b.Status() == third {
				return 1
			}
			return -1
		case fifth:
			if b.Status() == first || b.Status() == second || b.Status() == third || b.Status() == fourth {
				return 1
			}
			return -1
		default:
			return 0
		}
	}
}
