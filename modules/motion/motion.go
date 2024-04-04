package motion

import (
	"github.com/kjbreil/goscript/pkg/module"
)

var key = "motion"

type Motion struct {
	module.Base
	MotionLights map[string]motionLight
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
