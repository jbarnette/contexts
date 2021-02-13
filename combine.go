// Package contexts provides a helper to combine multiple instances of context.Context.
package contexts

import (
	"context"
	"sync"
	"time"
)

type combined struct {
	contexts []context.Context

	mu   sync.Mutex
	done chan struct{}
	err  error
}

// Combine returns a single context combining the provided contexts. The returned
// context's Done channel is closed the first time one of the provided contexts is Done.
// The returned context's Deadline is the oldest Deadline of any of the provided
// contexts. Values are also combined: The returned context's Value method calls Value on
// each provided context and returns the first non-nil result.
func Combine(contexts ...context.Context) context.Context {
	c := combined{
		contexts: contexts,
		done:     make(chan struct{}),
	}

	for _, ctx := range contexts {
		go c.wait(ctx)
	}

	return &c
}

func (c *combined) Deadline() (deadline time.Time, ok bool) {
	for _, ctx := range c.contexts {
		if d, has := ctx.Deadline(); has {
			if deadline.IsZero() || d.Before(deadline) {
				deadline, ok = d, true
			}
		}
	}

	return
}

func (c *combined) Done() <-chan struct{} {
	return c.done
}

func (c *combined) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.err
}

func (c *combined) Value(key interface{}) interface{} {
	for _, ctx := range c.contexts {
		if v := ctx.Value(key); v != nil {
			return v
		}
	}

	return nil
}

func (c *combined) wait(ctx context.Context) {
	select {
	case <-c.done:
		return
	case <-ctx.Done():
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.err != nil {
		return
	}

	c.err = ctx.Err()
	close(c.done)
}
