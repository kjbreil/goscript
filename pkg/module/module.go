package module

import (
	"context"
)

type Module interface {
	// Init brings in data streams that might be needed and returns the triggers the module provides
	Init(ctx context.Context, requests *Requests) error
	Responses(rsp Response)
	Requests() *Requests
	Run() error
	Update() error
	Name() string
	Close() error

	mustImplementBase()
}

type Base struct {
	Ctx      context.Context
	requests *Requests
}

func (m *Base) mustImplementBase() {}

func (m *Base) Init(ctx context.Context, requests *Requests) error {
	m.AssignBase(ctx, requests)

	return nil
}

func (m *Base) AssignBase(ctx context.Context, requests *Requests) {
	m.Ctx = ctx
	m.requests = requests
}

func (m *Base) Responses(_ Response) {
}
func (m *Base) Requests() *Requests {
	return m.requests
}
