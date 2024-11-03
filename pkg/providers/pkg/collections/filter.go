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

package collections

type Filter[T any] func(T) bool

func Slice[T any](filter Filter[T]) ListOption[T] {
	return func(ts []T) []T {
		var filtered []T
		for _, t := range ts {
			if filter(t) {
				filtered = append(filtered, t)
			}
		}
		return filtered
	}
}

func And[T any](filters ...Filter[T]) Filter[T] {
	return func(t T) bool {
		for _, filter := range filters {
			if !filter(t) {
				return false
			}
		}
		return true
	}
}

func Or[T any](filters ...Filter[T]) Filter[T] {
	return func(t T) bool {
		for _, filter := range filters {
			if filter(t) {
				return true
			}
		}
		return false
	}
}

func Not[T any](filter Filter[T]) Filter[T] {
	return func(t T) bool {
		return !filter(t)
	}
}
