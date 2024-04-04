package motion

import (
	"github.com/go-logr/logr"
	"github.com/kjbreil/goscript/pkg/core"
	"github.com/kjbreil/goscript/pkg/trigger"
)

type motionLight struct {
	logger logr.Logger

	triggers  []*trigger.Trigger
	service   core.ServiceChan
	TurnOn    bool
	Timeout   int
	BlockIfOn []string
	Detectors []string
	Entities  []string
}
