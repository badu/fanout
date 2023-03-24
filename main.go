package fanout

import (
	"sync"
)

type Fanner[T any] struct {
	mu     sync.Mutex
	chans  []chan T
	closed bool
}

func New[T any]() *Fanner[T] {
	result := Fanner[T]{
		chans: make([]chan T, 0),
	}
	return &result
}

func (c *Fanner[T]) Sub() <-chan T {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make(chan T)
	c.chans = append(c.chans, result)

	return result
}

func (c *Fanner[T]) Pub(payload T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed { // channel was not closed
		for _, ch := range c.chans {
			ch <- payload
		}
	}
}

func (c *Fanner[T]) Cancel(ch <-chan T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return true
	}

	found := false
	for i, ci := range c.chans {
		if ci == ch {
			c.chans[i] = c.chans[len(c.chans)-1]
			c.chans[len(c.chans)-1] = nil
			c.chans = c.chans[:len(c.chans)-1]
			found = true
			break
		}
	}
	return found
}

func (c *Fanner[T]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		c.chans = nil
		c.closed = true
	}
}
