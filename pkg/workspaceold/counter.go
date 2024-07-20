package workspace

import "sync/atomic"

type AtomicCounter struct {
	number uint64
}

func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{0}
}

func (c *AtomicCounter) Increment() {
	for {
		v := c.Read()
		if atomic.CompareAndSwapUint64(&c.number, v, v+1) {
			return
		}
	}
}

func (c *AtomicCounter) Decrement() {
	for {
		v := c.Read()
		if atomic.CompareAndSwapUint64(&c.number, v, v-1) {
			return
		}
	}
}

func (c *AtomicCounter) Read() uint64 {
	return atomic.LoadUint64(&c.number)
}
