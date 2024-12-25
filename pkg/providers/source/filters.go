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

package source

import (
	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
)

var And = collections.And[Source]
var Or = collections.Or[Source]

func WithFilter(filter collections.Filter[Source]) collections.ListOption[Source] {
	return func(tasks []Source) []Source {
		return collections.Slice(filter)(tasks)
	}
}

// Formats are OR'd together because a source can only have one status.
func FilterByFormat(status ...Format) collections.Filter[Source] {
	return func(source Source) bool {
		for _, s := range status {
			if source.Format == s {
				return true
			}
		}
		return false
	}
}
