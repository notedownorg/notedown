// Copyright 2025 Notedown Authors
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

package parser

// Visitor defines the visitor pattern interface for traversing the document tree
type Visitor interface {
	Visit(node Node) error
}

// Walker provides utilities for traversing document trees
type Walker struct {
	visitor Visitor
}

// NewWalker creates a new tree walker
func NewWalker(visitor Visitor) *Walker {
	return &Walker{visitor: visitor}
}

// Walk traverses the tree starting from the given node
func (w *Walker) Walk(node Node) error {
	if err := node.Accept(w.visitor); err != nil {
		return err
	}

	for _, child := range node.Children() {
		if err := w.Walk(child); err != nil {
			return err
		}
	}

	return nil
}

// WalkFunc is a convenience type for creating visitors from functions
type WalkFunc func(node Node) error

// Visit implements the Visitor interface
func (f WalkFunc) Visit(node Node) error {
	return f(node)
}
