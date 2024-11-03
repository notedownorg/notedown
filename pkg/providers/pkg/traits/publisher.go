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

package traits

type Publisher[Event any] struct {
	subscribers []chan Event
	Events      chan Event
}

func NewPublisher[Event any]() *Publisher[Event] {
	p := &Publisher[Event]{
		subscribers: make([]chan Event, 0),
		Events:      make(chan Event),
	}
	go p.start()
	return p
}

// Subscribe adds a new subscriber to the list of subscribers
func (p *Publisher[Event]) Subscribe(ch chan Event) int {
	p.subscribers = append(p.subscribers, ch)
	index := len(p.subscribers) - 1
	return index
}

// Unsubscribe removes a subscriber from the list of subscribers
func (p *Publisher[Event]) Unsubscribe(index int) {
	p.subscribers = append(p.subscribers[:index], p.subscribers[index+1:]...)
}

func (s *Publisher[Event]) start() {
	for event := range s.Events {
		for _, subscriber := range s.subscribers {
			go func(s chan Event, e Event) { s <- e }(subscriber, event)
		}
	}
}
