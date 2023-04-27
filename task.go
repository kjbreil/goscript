package goscript

import (
	"context"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/kjbreil/hass-ws/model"
	"sync"
	"time"
)

// Task is used within a TriggerFunc to give information about the task.
// Message is the message that triggered the task.
// States is all the States defined when the trigger was created. States gets updated each time the methods are run.
// Task contains 3 methods: Sleep, WaitUntil and While to help processing within a function and handle being able to
// properly kill the task externally.
type Task struct {
	Message     *model.Message
	MqttMessage mqtt.Message
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
	running     *bool
}

// TaskFunc is used to include a task object in MQTT command functions.
type TaskFunc func(t *Task)

// TaskMQTT wraps a trigger and TaskFunc setting up and passing the task through
func (gs *GoScript) TaskMQTT(tr *Trigger) func(message mqtt.Message, client mqtt.Client) {
	// setup the trigger
	tr = setupTrigger(tr)

	return func(message mqtt.Message, client mqtt.Client) {
		task := gs.newTask(tr, nil)
		task.MqttMessage = message
		gs.taskToRun.add(task)
	}
}

// Sleep waits for the timeout to occur and panics if the context is cancelled
// The panic is caught by a recover
func (t *Task) Sleep(timeout time.Duration) {
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		return
		//t.States = t.gs.GetStates(t.states)
	case <-t.ctx.Done():
		panic(fmt.Sprintf("task context cancelled for %s", t.uuid))
	}
}

// Context return the current tasks context
func (t *Task) Context() context.Context {
	return t.ctx
}

// UUID return the current tasks uuid
func (t *Task) UUID() uuid.UUID {
	return t.uuid
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
			//t.States = t.gs.GetStates(t.states)
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
		if eState, ok := t.States.Get(entityID); ok {
			if Evaluates(States{
				s: map[string]*State{entityID: eState},
				m: &sync.Mutex{},
			}, eval) {
				whileFunc()
			} else {
				return
			}
		} else {
			return
		}
	}
}

func (t *Task) Cancelled() bool {
	return errors.Is(t.ctx.Err(), context.Canceled)
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

func (gs *GoScript) runTask(t *Task) {
	for *t.running {
		timer := time.NewTimer(100)
		select {
		case <-timer.C:
		case <-t.ctx.Done():
			gs.logger.Info(fmt.Sprintf("task %s exited awaiting to run", t.uuid))
			return
		}
	}

	defer func() {
		*t.running = false
		t.cancel()
		if r := recover(); r != nil {
			gs.logger.Info(fmt.Sprintf("task exited: %v", r))
		}
	}()

	*t.running = true

	go gs.taskWaitRequest(t)
	t.f(t)
}

func (gs *GoScript) newTask(tr *Trigger, message *model.Message) *Task {
	task := &Task{
		Message:     message,
		States:      gs.states.SubSet(tr.States),
		ServiceChan: gs.ServiceChan,
		states:      tr.States,
		f:           tr.Func,
		waitRequest: make(chan *Trigger),
		waitDone:    make(chan bool),
	}

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
		task.running = new(bool)
		task.ctx, task.cancel = context.WithCancel(gs.ctx)
		task.uuid = uuid.New()
	}

	return task
}

func (gs *GoScript) makeUniqueTask(tr *Trigger, task *Task) (*Task, bool) {
	// KillMe checks if the task is running and exits rather than kill off the other task
	if tr.Unique.KillMe {
		if *tr.Unique.running {
			gs.logger.Info(fmt.Sprintf("task %s tried to start but other task running and KillMe is true", tr.uuid))
			return nil, true
		}
	}

	// non wait tasks (default) cancel the current context
	if !tr.Unique.Wait {
		tr.Unique.cancel()
		tr.Unique.ctx, tr.Unique.cancel = context.WithCancel(context.Background())
		task.ctx, task.cancel = tr.Unique.ctx, tr.Unique.cancel
	} else {
		task.ctx, task.cancel = context.WithCancel(tr.Unique.ctx)
		task.uuid = uuid.New()
	}

	task.running = tr.Unique.running

	if tr.Unique.UUID != nil {
		task.uuid = *tr.Unique.UUID
	} else {
		task.uuid = tr.uuid
	}
	return nil, false
}
