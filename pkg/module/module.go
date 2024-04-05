package module

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/service"
	"github.com/kjbreil/goscript/pkg/state"
	"github.com/kjbreil/goscript/pkg/trigger"
)

type Module interface {
	// Init brings in data streams that might be needed and returns the triggers the module provides
	Init(ctx context.Context,
		logger logr.Logger,
		sChan service.Chan,
		runner *trigger.Runner,
		states *state.States,
		modules *Modules) error
	Triggers() trigger.Triggers
	Devices() device.Devices
	Update() error
	Name() string
	Close() error

	mustImplementBase()
}

type Base struct {
	Ctx         context.Context
	ServiceChan service.Chan
	Runner      *trigger.Runner
	Logger      logr.Logger
	States      *state.States
	Modules     *Modules
}

func (m *Base) mustImplementBase() {}

func (m *Base) Init(ctx context.Context,
	logger logr.Logger,
	sChan service.Chan,
	runner *trigger.Runner,
	states *state.States,
	modules *Modules) error {
	m.AssignBase(ctx, logger, sChan, runner, states, modules)

	return nil
}

func (m *Base) AssignBase(ctx context.Context, logger logr.Logger, sChan service.Chan, runner *trigger.Runner, states *state.States, modules *Modules) {
	m.Ctx = ctx
	m.Logger = logger
	m.ServiceChan = sChan
	m.Runner = runner
	m.States = states
	m.Modules = modules
}
