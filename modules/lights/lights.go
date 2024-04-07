package lights

import (
	"github.com/go-logr/logr"
	"github.com/kjbreil/goscript/modules/circadian"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/goscript/pkg/service"
	"github.com/kjbreil/goscript/pkg/trigger"
)

var key = "lights"

type Lights struct {
	MotionLights map[string]motionLight
	service      service.Chan
	logger       logr.Logger
	cir          circadian.Circadian
	module.Base
}

func (l *Lights) Run() error {
	module.SendTriggers(l, l.motion())

	l.Requests().Chan() <- module.Request{
		To:      "goscript",
		From:    key,
		Trigger: nil,
		Device:  nil,
		Module:  "virtual",
	}

	module.SendInfo(l, "test message")

	return nil
}
func (l *Lights) Responses(rsp module.Response) {

}

func (l *Lights) Close() error {
	return nil
}

func (l *Lights) Update() error {

	return nil
}
func (l *Lights) Name() string {
	return key
}

func (l *Lights) Triggers() trigger.Triggers {
	var triggers trigger.Triggers
	// Setup the motion lights
	triggers = append(triggers, l.motion()...)
	return triggers
}

func (l *Lights) Devices() device.Devices {
	return nil
}
