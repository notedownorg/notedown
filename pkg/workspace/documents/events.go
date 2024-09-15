package documents

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
