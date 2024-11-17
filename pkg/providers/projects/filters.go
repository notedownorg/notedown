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

package projects

import (
	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
)

var And = collections.And[Project]
var Or = collections.Or[Project]

func WithFilter(filter collections.Filter[Project]) collections.ListOption[Project] {
	return func(tasks []Project) []Project {
		return collections.Slice(filter)(tasks)
	}
}

// Statuses are OR'd together because a project can only have one status.
func FilterByStatus(status ...Status) collections.Filter[Project] {
	return func(project Project) bool {
		for _, s := range status {
			if project.Status() == s {
				return true
			}
		}
		return false
	}
}
