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

type Event struct {
	Op Operation
}

type Operation uint32

const (
	// Signal that one or more tasks have been loaded
	Load Operation = iota

	// Signal that one or more tasks have been changed
	Change

	// Signal that one or more tasks have been deleted
	Delete
)

func (c *Client) Subscribe(ch chan Event) int {
	c.subscribers = append(c.subscribers, ch)
	index := len(c.subscribers) - 1
	return index
}

func (c *Client) Unsubscribe(index int) {
	c.subscribers = append(c.subscribers[:index], c.subscribers[index+1:]...)
}

func (c *Client) eventDispatcher() {
	for event := range c.events {
		for _, subscriber := range c.subscribers {
			go func(s chan Event, e Event) { s <- e }(subscriber, event)
		}
	}
}
