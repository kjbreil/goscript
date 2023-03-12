package goscript

import (
	"context"
	"github.com/kjbreil/hass-ws/model"
	"time"
)

type Task struct {
	Message *model.Message
	States  States
	// task context
	ctx    context.Context
	cancel context.CancelFunc
	//
	states      []string
	f           TriggerFunc
	waitRequest chan *Trigger
	waitDone    chan bool
	gs          *GoScript // TODO: Figure out how to not need gs in each task, needed for getstates on sleep right now
}

func (t *Task) Sleep(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		t.States = t.gs.GetStates(t.states)
		return true
	case <-t.ctx.Done():
		return false
	}
}

func (t *Task) WaitUntil(entityId string, eval []string, timeout time.Duration) bool {

	t.waitRequest <- &Trigger{
		Triggers: []string{entityId},
		Eval:     eval,
	}
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		select {
		case <-t.waitDone:
			t.States = t.gs.GetStates(t.states)
			return true
		case <-timer.C:
			t.cancel()
			return false
		case <-t.ctx.Done():
			return false
		}
	} else {
		select {
		case <-t.waitDone:
			t.States = t.gs.GetStates(t.states)
			return true
		case <-t.ctx.Done():
			return false
		}
	}

}

func (gs *GoScript) taskWaitRequest(t *Task) {
	var trigger *Trigger
	for {
		select {
		case trigger = <-t.waitRequest:
			// TODO: Validate entityid is valid
			trigger.Func = func(it *Task) {
				t.waitDone <- true
				gs.RemoveTrigger(trigger)
			}
			gs.AddTrigger(trigger)
		case <-t.ctx.Done():
			if trigger != nil {
				gs.RemoveTrigger(trigger)
			}
			return
		}
	}

}
