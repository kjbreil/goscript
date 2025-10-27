package trigger

import (
	"context"
	"fmt"

	"github.com/kjbreil/hass-ws/model"
)

// AddTrigger adds a trigger to the trigger map. There is no validation of a
func (r *Runner) AddTrigger(tr *Trigger) {
	tr = SetupTrigger(tr)
	r.triggerMu.Lock()
	defer r.triggerMu.Unlock()

	// for each entity add to the triggers map
	for _, et := range tr.Triggers {
		r.triggers[et] = append(r.triggers[et], tr)
	}

	// for each domain add to the domain trigger map
	for _, edt := range tr.DomainTrigger {
		r.domainTrigger[edt] = append(r.domainTrigger[edt], tr)
	}

	// for each periodic add to the periodic map
	// cron time is an array of triggers so multiple triggers can have same cron schedule
	for _, ep := range tr.Periodic {
		// check if the period is blank then run immediately and don't add to the map
		if ep == "" {
			task := r.NewTask(tr, nil)
			r.taskToRun.Add(task)
			continue
		}
		r.periodic[ep] = append(r.periodic[ep], tr)
	}

	for _, es := range tr.Services {
		r.serviceTriggers[es] = append(r.serviceTriggers[es], tr)
	}
}

// RemoveTrigger can be used to remove a trigger while program is running.
func (r *Runner) RemoveTrigger(t *Trigger) {
	r.triggerMu.Lock()
	defer r.triggerMu.Unlock()

	for _, et := range t.Triggers {
		for i, te := range r.triggers[et] {
			if te.UUID() == t.UUID() {
				r.triggers[et] = append(r.triggers[et][:i], r.triggers[et][i+1:]...)
				break
			}
		}
	}
}

// AddTriggers helper function to add multiple triggers
func (r *Runner) AddTriggers(triggers ...*Trigger) {
	for _, t := range triggers {
		r.AddTrigger(t)
	}
}

func (r *Runner) RunTriggers(message model.Message) {
	r.triggerMu.RLock()

	// Create copies of trigger slices while holding the lock
	var entityTriggers []*Trigger
	if tr, ok := r.triggers[message.DomainEntity()]; ok {
		entityTriggers = make([]*Trigger, len(tr))
		copy(entityTriggers, tr)
	}

	var domainTriggers []*Trigger
	if tr, ok := r.domainTrigger[message.Domain()]; ok {
		domainTriggers = make([]*Trigger, len(tr))
		copy(domainTriggers, tr)
	}

	r.triggerMu.RUnlock()

	// Process triggers without holding the lock
	for _, trr := range entityTriggers {
		r.triggerDomainEntity(&message, trr)
	}

	for _, trr := range domainTriggers {
		r.triggerDomain(&message, trr)
	}
}

func (r *Runner) triggerDomainEntity(message *model.Message, tr *Trigger) {
	passed := tr.Evaluate(message)
	if passed {
		task := r.NewTask(tr, message)
		r.taskToRun.Add(task)
	}
}

func (r *Runner) triggerDomain(message *model.Message, tr *Trigger) {
	passed := tr.Evaluate(message)
	if passed {
		task := r.NewTask(tr, message)
		r.taskToRun.Add(task)
	}
}

func (r *Runner) RunServiceTriggers(message model.Message) {
	if message.Event == nil || message.Event.Data == nil || message.Event.Data.ServiceData == nil {
		return
	}

	r.triggerMu.RLock()

	// Create copies of all service triggers for relevant entities while holding the lock
	allTriggers := make(map[string][]*Trigger)
	for _, entity := range message.Event.Data.ServiceData.EntityId {
		if tr, ok := r.serviceTriggers[entity]; ok {
			triggers := make([]*Trigger, len(tr))
			copy(triggers, tr)
			allTriggers[entity] = triggers
		}
	}

	r.triggerMu.RUnlock()

	// Process triggers without holding the lock
	for _, triggers := range allTriggers {
		for _, trigger := range triggers {
			r.triggerService(&message, trigger)
		}
	}
}

func (r *Runner) triggerService(message *model.Message, trigger *Trigger) {
	task := r.NewTask(trigger, message)
	r.taskToRun.Add(task)
}

func (r *Runner) NewTask(tr *Trigger, message *model.Message) *Task {
	task := New(tr)
	task.Message = message
	task.States = r.states.SubSet(tr.States)
	task.ServiceChan = r.sChan

	domainStates := r.states.GetDomainStates(tr.DomainTrigger)
	task.States.Combine(domainStates)

	domainStates = r.states.GetDomainStates(tr.DomainStates)
	task.States.Combine(domainStates)

	if tr.Unique != nil {
		t, done := r.MakeUniqueTask(tr, task)
		if done {
			return t
		}
	} else {
		task.SetRunningFalse()
		task.NewCtx(r.ctx)
		task.NewUUID()
	}

	return task
}

func (r *Runner) MakeUniqueTask(tr *Trigger, task *Task) (*Task, bool) {
	// KillMe checks if the task is running and exits rather than kill off the other task
	if tr.Unique.KillMe {
		if *tr.Unique.Running() {
			r.logger.Info(fmt.Sprintf("task %s tried to start but other task running and KillMe is true", tr.UUID()))
			return nil, true
		}
	}
	// else {
	//
	// }

	// non wait tasks (default) cancel the current context
	if !tr.Unique.Wait {
		task.CopyCtx(tr.Unique.CancelNew(context.Background()))
	} else {
		task.NewCtx(r.ctx)
		task.NewUUID()
	}

	task.SetRunning(tr.Unique.Running())

	// TODO: This needs to be moved up, right now Unique.UUID only working for Wait = true
	if tr.Unique.UUID != nil {
		task.NewUUID()
		task.SetRunning(r.triggerRunning.Get(*tr.Unique.UUID))
	} else {
		// task.NewUUID()
		task.uuid = tr.uuid
	}
	return nil, false
}
