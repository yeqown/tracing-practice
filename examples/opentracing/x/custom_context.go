package x

import (
	"context"
	"time"
)

type customContext struct {
	// 当需要被设置 context 包的方法时，用context.Context 来传递和接收
	// 但是 Value 方法可以直接使用，因为也会被透传
	context.Context

	Field1 string
}

func newContext() *customContext {
	return &customContext{
		Context: context.Background(),
	}
}

func (c *customContext) WithTimeout(timeout time.Duration) *customContext {
	// copy c expect of context.Context
	c2 := &customContext{
		Field1: c.Field1,
	}
	c2.Context, _ = context.WithTimeout(c.Context, timeout)
	return c2
}

func (c customContext) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Deadline()
}

func (c customContext) Done() <-chan struct{} {
	return c.Context.Done()
}

func (c customContext) Err() error {
	return c.Context.Err()
}

func (c customContext) Value(key interface{}) interface{} {
	return c.Context.Value(key)
}
