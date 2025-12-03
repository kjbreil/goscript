package trigger

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/adhocore/gronx"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kjbreil/goscript/pkg/service"
	"github.com/kjbreil/goscript/pkg/state"
)

const (
	minutesPerHour       = 60
	taskRunCheckInterval = 100 * time.Millisecond
)

type Runner struct {
	// maps holding state based triggers
	periodic        map[string]Triggers
	nextPeriodic    time.Time
	triggers        map[string]Triggers
	domainTrigger   map[string]Triggers
	serviceTriggers map[string]Triggers

	triggerRunning Running

	taskToRun TaskMap
	states    *state.States
	sChan     service.Chan
	logger    *slog.Logger
	ctx       context.Context
	triggerMu sync.RWMutex
}

func NewRunner(ctx context.Context, states *state.States, sChan service.Chan, logger *slog.Logger) *Runner {
	//nolint:exhaustruct // nextPeriodic and triggerMu initialized with zero values
	return &Runner{
		triggers:        make(map[string]Triggers),
		domainTrigger:   make(map[string]Triggers),
		serviceTriggers: make(map[string]Triggers),
		periodic:        make(map[string]Triggers),
		taskToRun:       NewTaskMap(),
		triggerRunning:  NewRunning(),
		states:          states,
		sChan:           sChan,
		logger:          logger,
		ctx:             ctx,
	}
}

func (r *Runner) AddTask(t *Task) {
	r.taskToRun.Add(t)
}
func (r *Runner) TaskToRun() []*Task {
	return r.taskToRun.ToRun()
}

func (r *Runner) RunPeriodic() {
	var err error
	// TODO: Validate Periodic slice
	// run zero length immediate periodics and delete from periodic list
	// TODO: Validate that this still is used since now the zero length periodics should be added to the taskToRun immediately
	for _, triggers := range r.periodic {
		for _, t := range triggers {
			pLen := len(t.Periodic)
			for i := 0; i < pLen; i++ {
				if len(t.Periodic[i]) == 0 {
					task := r.NewTask(t, nil)
					r.taskToRun.Add(task)

					t.Periodic = append(t.Periodic[:i], t.Periodic[i+1:]...)
					i--
					pLen--
				}
			}
		}
	}
	delete(r.periodic, "")

	// setup the next fire time for all triggers
	r.nextPeriodic, err = fillNextTime(r.periodic)
	if err != nil {
		r.logger.Error("failed to get next time", "error", err.Error())
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				if now.After(r.nextPeriodic) {
					go r.shouldRunTrigger()
				}
			case <-r.ctx.Done():
				return
			}
		}
	}()
}

func (r *Runner) shouldRunTrigger() {
	r.nextPeriodic = time.Now().Add(minutesPerHour * time.Minute)
	for _, triggers := range r.periodic {
		for _, t := range triggers {
			if t.GetNextTime() == nil {
				r.logger.Info("next time not set")
				_, err := t.NextTime(time.Now())
				if err != nil {
					r.logger.Error("setting next time failed", "error", err.Error())
					continue
				}
			}
			if time.Now().After(*t.GetNextTime()) {
				task := r.NewTask(t, nil)
				r.taskToRun.Add(task)

				_, err := t.NextTime(time.Now())
				if err != nil {
					r.logger.Error("setting next time failed", "error", err.Error())
					continue
				}
			}
			if t.GetNextTime().Before(r.nextPeriodic) {
				r.nextPeriodic = *t.GetNextTime()
			}
		}
	}
}

func (r *Runner) runGronJob(gron *gronx.Gronx, start bool) {
	for expr, triggers := range r.periodic {
		var err error
		var due bool
		if len(expr) == 0 {
			if start {
				due = true
			}
		} else {
			due, err = gron.IsDue(expr)
			if err != nil {
				r.logger.Error("gron job IsDue failed", "error", err.Error())
				continue
			}
		}
		for _, t := range triggers {
			if due {
				task := r.NewTask(t, nil)
				r.taskToRun.Add(task)
			}
		}
	}
}

func (r *Runner) RunTask(t *Task) {
	r.logger.Debug(fmt.Sprintf("task started: %s", t.UUID()), "messageEntity", t.Message.DomainEntity())

	for t.Running() {
		timer := time.NewTimer(taskRunCheckInterval)
		select {
		case <-timer.C:
			timer.Stop()
		case <-t.CtxDone():
			timer.Stop()
			r.logger.Debug(fmt.Sprintf("task %s exited awaiting to run", t.UUID()))
			return
		}
	}

	defer func() {
		t.SetRunningFalse()
		t.Cancel()
		if re := recover(); re != nil {
			r.logger.Debug(fmt.Sprintf("task exited: %v", re))
		}
	}()

	t.SetRunningTrue()

	go r.taskWaitRequest(t)
	t.F()(t)
}

func (r *Runner) taskWaitRequest(t *Task) {
	var tr *Trigger
	for {
		select {
		case tr = <-t.WaitRequest():
			// TODO: Validate entityid is valid
			tr.Func = func(it *Task) {
				t.WaitDone() <- true
				r.RemoveTrigger(tr)
			}
			r.AddTrigger(tr)
		case <-t.CtxDone():
			if tr != nil {
				r.RemoveTrigger(tr)
			}
			return
		}
	}
}

func fillNextTime(periodics map[string]Triggers) (time.Time, error) {
	next := time.Now().Add(minutesPerHour * time.Minute)
	if len(periodics) == 0 {
		return time.Now().Add(1 * time.Second), nil
	}
	for _, triggers := range periodics {
		for _, t := range triggers {
			nt, err := t.NextTime(time.Now())
			if err != nil {
				return next, fmt.Errorf("failed to get NextTime for task %s: %w", t.UUID(), err)
			}
			if nt != nil && nt.Before(next) {
				next = *nt
			}
		}
	}
	return next, nil
}

func (r *Runner) TaskMQTT(tr *Trigger) func(message mqtt.Message, client mqtt.Client) {
	// setup the trigger
	tr = SetupTrigger(tr)

	return func(message mqtt.Message, client mqtt.Client) {
		task := r.NewTask(tr, nil)
		task.MqttMessage = message
		r.AddTask(task)
	}
}
