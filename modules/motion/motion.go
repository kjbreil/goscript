package motion

import (
	"github.com/kjbreil/goscript/pkg/core"
	"github.com/kjbreil/goscript/pkg/module"
)

var key = "motion"

type Motion struct {
	module.Base
}

func New(gs *core.GoScript) *Motion {

	m, err := gs.GetModule(key)

	return &m
}

func (m *Motion) Update() error {
	// TODO implement me
	panic("implement me")
}

func (m *Motion) Name() string {
	return key
}

func (m *Motion) Close() error {
	// TODO implement me
	panic("implement me")
}
