package goscript

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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
	uuid        uuid.UUID
	waitRequest chan *Trigger
	waitDone    chan bool
	gs          *GoScript // TODO: Figure out how to not need gs in each task, needed for getstates on sleep right now
}

// Sleep waits for the timeout to occur and panics if the context is cancelled
func (t *Task) Sleep(timeout time.Duration) {
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		t.States = t.gs.GetStates(t.states)
	case <-t.ctx.Done():
		panic(fmt.Sprintf("task context cancelled for %s", t.uuid))
	}
}

// WaitUntil waits
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
			panic(fmt.Sprintf("task context cancelled for %s", t.uuid))
		}
	} else {
		select {
		case <-t.waitDone:
			t.States = t.gs.GetStates(t.states)
			return true
		case <-t.ctx.Done():
			panic(fmt.Sprintf("task context cancelled for %s", t.uuid))
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

func (t *Task) run() {
	defer func() {
		if r := recover(); r != nil {
			t.gs.logger.Info(fmt.Sprintf("task exited: %v", r))
		}
	}()
	go t.gs.taskWaitRequest(t)

	t.f(t)
}

func (gs *GoScript) newTask(tr *Trigger, message *model.Message) *Task {
	task := &Task{
		Message:     message,
		States:      gs.GetStates(tr.States),
		gs:          gs,
		states:      tr.States,
		f:           tr.Func,
		waitRequest: make(chan *Trigger),
		waitDone:    make(chan bool),
	}
	if tr.Unique != nil {
		tr.Unique.cancel()
		tr.Unique.ctx, tr.Unique.cancel = context.WithCancel(context.Background())
		task.ctx, task.cancel = tr.Unique.ctx, tr.Unique.cancel
		task.uuid = tr.uuid
	} else {
		task.ctx, task.cancel = context.WithCancel(context.Background())
		newUUID := uuid.New()
		task.uuid = newUUID
	}

	return task
}
