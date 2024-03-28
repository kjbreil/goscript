package core

import (
	"context"
	"fmt"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/model"
	"time"
)

func (gs *GoScript) taskWaitRequest(t *trigger.Task) {
	var tr *trigger.Trigger
	for {
		select {
		case tr = <-t.WaitRequest():
			// TODO: Validate entityid is valid
			tr.Func = func(it *trigger.Task) {
				t.WaitDone() <- true
				gs.RemoveTrigger(tr)
			}
			gs.AddTrigger(tr)
		case <-t.CtxDone():
			if tr != nil {
				gs.RemoveTrigger(tr)
			}
			return
		}
	}
}

func (gs *GoScript) runTask(t *trigger.Task) {
	for t.Running() {
		timer := time.NewTimer(100)
		select {
		case <-timer.C:
		case <-t.CtxDone():
			gs.logger.Info(fmt.Sprintf("task %s exited awaiting to run", t.UUID()))
			return
		}
	}

	defer func() {
		t.SetRunning(false)
		t.Cancel()
		if r := recover(); r != nil {
			gs.logger.Info(fmt.Sprintf("task exited: %v", r))
		}
	}()

	t.SetRunning(true)

	go gs.taskWaitRequest(t)
	t.F()(t)
}

func (gs *GoScript) newTask(tr *trigger.Trigger, message *model.Message) *trigger.Task {
	task := trigger.New(tr)
	task.Message = message
	task.States = gs.states.SubSet(tr.States)
	task.ServiceChan = gs.ServiceChan

	domainStates := gs.GetDomainStates(tr.DomainTrigger)
	task.States.Combine(domainStates)

	domainStates = gs.GetDomainStates(tr.DomainStates)
	task.States.Combine(domainStates)

	if tr.Unique != nil {
		t, done := gs.makeUniqueTask(tr, task)
		if done {
			return t
		}
	} else {
		task.SetRunning(false)
		task.NewCtx(gs.ctx)
		task.NewUUID()
	}

	return task
}

func (gs *GoScript) makeUniqueTask(tr *trigger.Trigger, task *trigger.Task) (*trigger.Task, bool) {
	// KillMe checks if the task is running and exits rather than kill off the other task
	if tr.Unique.KillMe {
		if tr.Unique.Running() {
			gs.logger.Info(fmt.Sprintf("task %s tried to start but other task running and KillMe is true", tr.UUID()))
			return nil, true
		}
	}

	// non wait tasks (default) cancel the current context
	if !tr.Unique.Wait {
		task.CopyCtx(tr.Unique.CancelNew(context.Background()))
	} else {
		task.NewCtx(gs.ctx)
		task.NewUUID()
	}

	task.SetRunning(tr.Unique.Running())

	// TODO: This needs to be moved up, right now Unique.UUID only working for Wait = true
	if tr.Unique.UUID != nil {
		task.NewUUID()
		task.SetRunning(*gs.triggerRunning.Get(*tr.Unique.UUID))
	} else {
		task.NewUUID()
	}
	return nil, false
}
