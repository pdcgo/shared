package gorm_commenter

import (
	"context"
	"time"
)

type CommenterValue struct {
	Action      string
	Controller  string
	Framework   string
	Route       string
	Application string
	DbDriver    string

	ctx context.Context
}

// Deadline implements context.Context.
func (c *CommenterValue) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done implements context.Context.
func (c *CommenterValue) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err implements context.Context.
func (c *CommenterValue) Err() error {
	return c.ctx.Err()
}

// Value implements context.Context.
func (c *CommenterValue) Value(key any) any {
	switch key {
	case "action":
		return c.Action
	case "controller":
		return c.Controller
	case "framework":
		return c.Framework
	case "route":
		return c.Route
	case "application":
		return c.Application
	case "db driver":
		return c.DbDriver
	default:
		return c.ctx.Value(key)
	}
}

func NewCommenterContext(ctx context.Context, val CommenterValue) context.Context {
	val.ctx = ctx
	return &val
}
