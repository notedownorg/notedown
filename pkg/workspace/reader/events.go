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
	"log/slog"

	"github.com/notedownorg/notedown/pkg/workspace"
)

type Event struct {
	Op       Operation
	Key      string
	Document workspace.Document
}

type Operation uint32

const (
	// Signal that this document was present when the client was created or when the subscriber subscribed
	Load Operation = iota

	// Signal that this document has been updated or created
	Change

	// Signal that this document has been deleted
	Delete

	// Signal that the subscriber has received all existing documents present at the time of subscription
	SubscriberLoadComplete
)

type subscribeOptions func(*Client, chan Event)

// Load all existing documents as events to the new subscriber.
// Once all events have been sent, the a LoadComplete event is sent.
func WithInitialDocuments() subscribeOptions {
	return func(client *Client, sub chan Event) {
		go func(s chan Event) {
			for key, doc := range client.documents {
				s <- Event{Op: Load, Document: doc, Key: key}
			}
			s <- Event{Op: SubscriberLoadComplete}
		}(sub)
	}
}

func (c *Client) Subscribe(ch chan Event, opts ...subscribeOptions) int {
	c.subscribers = append(c.subscribers, ch)
	index := len(c.subscribers) - 1

	// Apply any subscribeOptions
	for _, opt := range opts {
		opt(c, ch)
	}
	return index
}

func (c *Client) Unsubscribe(index int) {
	c.subscribers = append(c.subscribers[:index], c.subscribers[index+1:]...)
}

func (c *Client) eventDispatcher() {
	for event := range c.events {
		for _, subscriber := range c.subscribers {
			go func(s chan Event, e Event) {
				// Recover from a panic if the subscriber has been closed
				// Likely this will only happen in tests but its theoretically possible in regular usage
				defer func() {
					if r := recover(); r != nil {
						slog.Warn("dropping event as subscriber has been closed", "path", e.Key)
					}
				}()
				s <- e
			}(subscriber, event)
		}
	}
}
