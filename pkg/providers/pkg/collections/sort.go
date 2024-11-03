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

import (
	"slices"
	"strings"
)

type Sorter[T any] func(a, b T) int

func Sort[T any](sorter Sorter[T]) ListOption[T] {
	return func(ts []T) []T {
		slices.SortFunc(ts, sorter)
		return ts
	}
}

// Fallthrough returns a sorter that falls through the given sorters until a non-zero result is found.
// If all sorters return 0, 0 is returned.
func Fallthrough[T any](sorters ...Sorter[T]) Sorter[T] {
	return func(a, b T) int {
		for _, sorter := range sorters {
			if result := sorter(a, b); result != 0 {
				return result
			}
		}
		return 0
	}
}

// FallthroughDeterministic returns a sorter that falls through the given sorters until a non-zero result is found.
// If all sorters return 0, the results are sorted alphabetically by the Name() value.
func FallthroughDeterministic[T nameable](sorters ...Sorter[T]) Sorter[T] {
	sorters = append(sorters, sortByAlphabetical[T]())
	return Fallthrough[T](sorters...)
}

type nameable interface {
	Name() string
}

func sortByAlphabetical[T nameable]() Sorter[T] {
	return func(a, b T) int {
		return strings.Compare(a.Name(), b.Name())
	}
}
