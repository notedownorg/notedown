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

package daily

import (
	"time"

	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
)

func WithFilters(filters ...collections.Filter[Daily]) collections.ListOption[Daily] {
	return func(tasks []Daily) []Daily {
		return collections.Slice[Daily](collections.And(filters...))(tasks)
	}
}

// Following Go's time package, after and before are inclusive (include equal to).
func FilterByDate(after *time.Time, before *time.Time) collections.Filter[Daily] {
	return func(d Daily) bool {
		if after != nil && d.date.Before(*after) {
			return false
		}
		if before != nil && d.date.After(*before) {
			return false
		}
		return true
	}
}
