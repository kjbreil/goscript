package goscript

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kjbreil/hass-ws/model"
	"sync"
	"time"
)

// Task is used within a TriggerFunc to give information about the task.
// Message is the message that triggered the task.
// States is all the states defined when the trigger was created. States gets updated each time the methods are run.
// Task contains 3 methods: Sleep, WaitUntil and While to help processing within a function and handle being able to
// properly kill the task externally.
type Task struct {
	Message     *model.Message
	States      States
	ServiceChan ServiceChan
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

type taskRun struct {
	tasks map[uuid.UUID]*Task
	m     *sync.Mutex
}

func (tr *taskRun) add(t *Task) {
	tr.m.Lock()
	defer tr.m.Unlock()
	tr.tasks[t.uuid] = t
}

// Sleep waits for the timeout to occur and panics if the context is cancelled
// The panic is caught by a recover
func (t *Task) Sleep(timeout time.Duration) {
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		t.States = t.gs.GetStates(t.states)
	case <-t.ctx.Done():
		panic(fmt.Sprintf("task context cancelled for %s", t.uuid))
	}
}

// WaitUntil waits until the eval equals true. Timeout of 0 means no timeout
// panics if the context is cancelled
func (t *Task) WaitUntil(entityID string, eval []string, timeout time.Duration) bool {

	t.waitRequest <- &Trigger{
		Triggers: []string{entityID},
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

// WhileFunc is the function that runs inside of a task.While on a continuous loop until the while evals false
type WhileFunc func()

// While runs a function until the eval does not evaluate true
// panics if the context is cancelled
// take care to use a sleep within the whileFunc
// best to keep the function inline so task.Sleep can be used
func (t *Task) While(entityID string, eval []string, whileFunc WhileFunc) {
	for {
		if t.ctx.Err() != nil {
			panic(fmt.Sprintf("task context cancelled for %s", t.uuid))
		}
		t.States = t.gs.GetStates(t.states)
		if eState, ok := t.States[entityID]; ok {
			if Evaluates(map[string]*State{entityID: eState}, eval) {
				whileFunc()
			} else {
				return
			}
		} else {
			return
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
		t.cancel()
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

	domainStates := gs.GetDomainStates(tr.DomainTrigger)
	for k, v := range domainStates {
		task.States[k] = v
		task.states = append(task.states, k)
	}

	domainStates = gs.GetDomainStates(tr.DomainStates)
	for k, v := range domainStates {
		task.States[k] = v
		task.states = append(task.states, k)
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
