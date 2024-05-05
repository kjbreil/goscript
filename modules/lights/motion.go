package lights

import (
	"github.com/go-logr/logr"
	"github.com/kjbreil/goscript/light"
	"github.com/kjbreil/goscript/modules/circadian"
	"github.com/kjbreil/goscript/pkg/eval"
	"github.com/kjbreil/goscript/pkg/service"
	"github.com/kjbreil/goscript/pkg/trigger"
	"time"
)

type motionLight struct {
	logger logr.Logger
	cir    *circadian.Circadian

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
		triggers = append(triggers, l.mlTrigger(ml))
	}

	return triggers
}

func (l *Lights) mlTrigger(ml *motionLight) *trigger.Trigger {
	return &trigger.Trigger{
		Triggers: ml.Detectors,

		Unique: &trigger.Unique{},
		States: append(ml.Entities, ml.BlockIfOn...),
		Eval:   eval.Eval(`state == "on"`),
		Func: func(t *trigger.Task) {
			turnOn := ml.TurnOn
			for _, e := range ml.BlockIfOn {
				for _, s := range t.States.Slice() {
					if e == s.DomainEntity && s.State == "on" {
						turnOn = false
					}
				}
			}
			if turnOn {
				if l.cir != nil {
					l.cir.TurnOn(t, ml.Entities...)
				} else {
					light.New().
						BrightnessPct(100).
						TurnOn(t, ml.Entities)
				}
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
			t.WaitUntil(t.Message.DomainEntity(), eval.Eval(`state == "off"`), 0)

			t.Sleep(time.Duration(ml.Timeout) * time.Second)

			light.New().TurnOff(t, ml.Entities)
		},
	}
}
