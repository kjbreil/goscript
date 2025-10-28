package lights

import (
	"github.com/kjbreil/goscript/modules/circadian"
	"github.com/kjbreil/goscript/pkg/control"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/goscript/pkg/trigger"
)

var key = "lights"

type Lights struct {
	MotionLights map[string]*motionLight
	cir          *circadian.Circadian
	module.Base
}

func (l *Lights) Run() error {
	module.SendTriggers(l, l.motion())

	l.Requests().Chan() <- control.Request{
		To:      "goscript",
		From:    key,
		Trigger: nil,
		Device:  nil,
		Module:  &circadian.Circadian{},
		Callback: func(rsp control.Response) error {
			l.cir = rsp.Module.(*circadian.Circadian)
			for _, ml := range l.MotionLights {
				ml.cir = l.cir
			}
			return nil
		},
	}

	module.SendInfo(l, "test message")

	return nil
}
func (l *Lights) Responses(rsp control.Response) {

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
