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

import (
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
)

type EventHandler func(reader.Event)

// WithInitialLoadWaiter waits for the initial load to complete before returning
func WithInitialLoadWaiter(watcher *Watcher) func(tick time.Duration) {
	return func(tick time.Duration) {
		for !watcher.InitialLoadComplete {
			time.Sleep(tick)
		}
	}
}

// A subcriber subscribes to events from the fileserver reader
// To use the trait you must provide actions to do on each event type
type Watcher struct {
	feed <-chan reader.Event

	onLoad   func(reader.Event)
	onChange func(reader.Event)
	onDelete func(reader.Event)

	InitialLoadComplete bool
}

func NewWatcher(feed <-chan reader.Event, onLoad, onChange, onDelete EventHandler) *Watcher {
	s := &Watcher{
		feed:                feed,
		onLoad:              onLoad,
		onChange:            onChange,
		onDelete:            onDelete,
		InitialLoadComplete: false,
	}
	go s.start()
	return s
}

func (s *Watcher) start() {
	for {
		select {
		case event := <-s.feed:
			switch event.Op {
			case reader.Delete:
				s.onDelete(event)
			case reader.Change:
				s.onChange(event)
			case reader.Load:
				s.onLoad(event)
			case reader.SubscriberLoadComplete:
				s.InitialLoadComplete = true
			}
		}
	}
}
