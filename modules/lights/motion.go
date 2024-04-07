package lights

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/kjbreil/goscript/light"
	"github.com/kjbreil/goscript/modules/circadian"
	"github.com/kjbreil/goscript/pkg/eval"
	"github.com/kjbreil/goscript/pkg/service"
	"github.com/kjbreil/goscript/pkg/trigger"
	"strings"
	"time"
)

type motionLight struct {
	logger logr.Logger
	cir    circadian.Circadian

	triggers  trigger.Triggers
	service   service.Chan
	TurnOn    bool
	Timeout   int
	BlockIfOn []string
	Detectors []string
	Entities  []string
}

func (l *Lights) motion() trigger.Triggers {
	var triggers trigger.Triggers

	for _, ml := range l.MotionLights {
		ml.logger = l.logger
		triggers = append(triggers, ml.trigger())
	}

	return triggers
}

func (l *motionLight) trigger() *trigger.Trigger {
	return &trigger.Trigger{
		Triggers: l.Detectors,

		Unique: &trigger.Unique{},
		States: append(l.Entities, l.BlockIfOn...),
		Eval:   eval.Eval(`state == "on"`),
		Func:   l.turnOnLights,
	}
}

func (l *motionLight) turnOnLights(t *trigger.Task) {
	turnOn := l.TurnOn
	l.logger.Info(fmt.Sprintf("Motion Detected: %s", t.Message.DomainEntity()))
	for _, e := range l.BlockIfOn {
		for _, s := range t.States.Slice() {
			if e == s.DomainEntity && s.State == "on" {
				turnOn = false
			}
		}
	}
	if turnOn {
		l.logger.Info(fmt.Sprintf("Turning Lights on: %s", strings.Join(l.Entities, ",")))
		// l.cir.TurnOn(t, l.Entities...)
		light.New().
			BrightnessPct(100).
			TurnOn(t, l.Entities)
	} else {
		allOff := true
		for _, s := range t.States.Slice() {
			if s.Domain == "light" && s.State == "on" {
				allOff = false
			}
		}
		if allOff {
			return
		}
	}
	l.logger.Info(fmt.Sprintf("Waiting for no motion on %s", t.Message.DomainEntity()))
	t.WaitUntil(t.Message.DomainEntity(), eval.Eval(`state == "off"`), 0)

	l.logger.Info(fmt.Sprintf("Sleeping for %d entity: %s", l.Timeout, t.Message.DomainEntity()))
	t.Sleep(time.Duration(l.Timeout) * time.Second)

	l.logger.Info(fmt.Sprintf("Turning off lights: %s", strings.Join(l.Entities, ",")))
	light.New().TurnOff(t, l.Entities)
	// l.circadian.TurnOff(l.Entities...)
}
