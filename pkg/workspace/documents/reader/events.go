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

import (
	"sync"
)

type Event struct {
	Op       Op
	Key      string
	Document Document
}

type Op uint32

const (
	Load   Op = iota
	Change    // Use change rather then create/update to promote idempotency
	Delete
)

type subscribeOptions func(*Client, chan Event)

// Load all existing documents as events to the new subscriber (events will be of Op Load).
// If a waitgroup is provided, it will be incremented by one for each document
// This way the caller can choose if/when/how to wait for all initial documents to be sent
func WithInitialDocuments(wg *sync.WaitGroup) subscribeOptions {
	return func(client *Client, sub chan Event) {
		for key, doc := range client.documents {
			if wg != nil {
				wg.Add(1)
			}
			go func(s chan Event, d Document, k string) {
				s <- Event{Op: Load, Document: d, Key: k}
			}(sub, doc, key)
		}
	}
}

func (c *Client) Subscribe(ch chan Event, opts ...subscribeOptions) {
	c.subscribers = append(c.subscribers, ch)

	// Apply any subscribeOptions
	for _, opt := range opts {
		opt(c, ch)
	}
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
