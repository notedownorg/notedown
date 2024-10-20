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

package reader

type Event struct {
	Op       Op
	Key      string
	Document Document
}

type Op uint32

const (
	// Use change rather then create/update to promote idempotency
	Change Op = iota
	Delete
)

func (c *Client) Subscribe() <-chan Event {
	sub := make(chan Event)
	c.subscribers = append(c.subscribers, sub)
	return sub
}

func (c *Client) eventDispatcher() {
	for {
		select {
		case ev := <-c.events:
			for _, sub := range c.subscribers {
				go func(s chan Event, e Event) { s <- e }(sub, ev)
			}
		}
	}
}
