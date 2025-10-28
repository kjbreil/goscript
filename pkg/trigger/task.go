package trigger

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/kjbreil/goscript/pkg/eval"
	"github.com/kjbreil/goscript/pkg/service"
	"github.com/kjbreil/goscript/pkg/state"
	"github.com/kjbreil/hass-ws/model"
)

// Task is used within a TriggerFunc to give information about the task.
// Message is the message that triggered the task.
// States is all the States defined when the trigger was created. States gets updated each time the methods are run.
// Task contains 3 methods: Sleep, WaitUntil and While to help processing within a function and handle being able to
// properly kill the task externally.
type Task struct {
	Message     *model.Message
	MqttMessage mqtt.Message
	States      state.States
	ServiceChan service.Chan
	// requests    *control.Requests

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
	runningMu   sync.RWMutex // Protects access to running pointer and its value
}

// TaskFunc is used to include a task object in MQTT command functions.
type TaskFunc func(t *Task)

func New(tr *Trigger) *Task {
	return &Task{
		states:      tr.States,
		f:           tr.Func,
		waitRequest: make(chan *Trigger, 1), // Buffered to prevent deadlock
		waitDone:    make(chan bool, 1),     // Buffered to prevent deadlock
	}
}

func (t *Task) NewCtx(ctx context.Context) {
	t.ctx, t.cancel = context.WithCancel(ctx)
}
func (t *Task) CopyCtx(ctx context.Context, cancel context.CancelFunc) {
	t.ctx, t.cancel = ctx, cancel
}

func (t *Task) NewUUID() {
	t.uuid = uuid.New()
}

func (t *Task) F() TriggerFunc {
	return t.f
}

func (t *Task) WaitRequest() chan *Trigger {
	return t.waitRequest
}

func (t *Task) WaitDone() chan bool {
	return t.waitDone
}
func (t *Task) Cancel() {
	t.cancel()
}

func (t *Task) Running() bool {
	t.runningMu.RLock()
	defer t.runningMu.RUnlock()
	if t.running == nil {
		return false
	}
	return *t.running
}

func (t *Task) SetRunning(running *bool) {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	t.running = running
}

func (t *Task) SetRunningTrue() {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	if t.running == nil {
		t.running = new(bool)
	}

	*t.running = true
}

func (t *Task) SetRunningFalse() {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	if t.running == nil {
		t.running = new(bool)
	}

	*t.running = false
}
func (t *Task) CtxDone() <-chan struct{} {
	return t.ctx.Done()
}

// Sleep waits for the timeout to occur and panics if the context is cancelled
// The panic is caught by a recover.
func (t *Task) Sleep(timeout time.Duration) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return
		// t.States = t.gs.GetStates(t.states)
	case <-t.ctx.Done():
		panic(fmt.Sprintf("task context cancelled for %s", t.uuid))
	}
}

// Context return the current tasks context.
func (t *Task) Context() context.Context {
	return t.ctx
}

// UUID return the current tasks uuid.
func (t *Task) UUID() uuid.UUID {
	return t.uuid
}

// WaitUntil waits until the eval equals true. Timeout of 0 means no timeout
// panics if the context is cancelled.
func (t *Task) WaitUntil(entityID string, eval []string, timeout time.Duration) bool {
	t.waitRequest <- &Trigger{
		Triggers: []string{entityID},
		Eval:     eval,
	}
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		select {
		case <-t.waitDone:
			// t.States = t.gs.GetStates(t.states)
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

// WhileFunc is the function that runs inside of a task.While on a continuous loop until the while evals false.
type WhileFunc func()

// While runs a function until the eval does not evaluate true
// panics if the context is cancelled
// take care to use a sleep within the whileFunc
// best to keep the function inline so task.Sleep can be used.
func (t *Task) While(entityID string, ev []string, whileFunc WhileFunc) {
	for {
		if t.ctx.Err() != nil {
			panic(fmt.Sprintf("task context cancelled for %s", t.uuid))
		}
		if eState, ok := t.States.Get(entityID); ok {
			if eval.Evaluates(state.NewSingleStates(entityID, eState), ev) {
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

func ptr[T any](v T) *T {
	return &v
}
