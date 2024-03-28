package core

import (
	"fmt"
	"github.com/adhocore/gronx"
	"github.com/kjbreil/goscript/pkg/trigger"
	"time"
)

func (gs *GoScript) runPeriodic() {
	var err error
	// TODO: Validate Periodic slice
	// run zero length immediate periodics and delete from periodic list
	for _, triggers := range gs.periodic {
		for _, t := range triggers {
			pLen := len(t.Periodic)
			for i := 0; i < pLen; i++ {
				if len(t.Periodic[i]) == 0 {
					task := gs.newTask(t, nil)
					gs.taskToRun.Add(task)

					t.Periodic = append(t.Periodic[:i], t.Periodic[i+1:]...)
					i--
					pLen--
				}
			}
		}
	}
	delete(gs.periodic, "")

	// setup the next fire time for all triggers
	gs.nextPeriodic, err = fillNextTime(gs.periodic)
	if err != nil {
		gs.Logger().Error(err, "NextTime")
	}

	//

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if time.Now().After(gs.nextPeriodic) {
				go gs.shouldRunTrigger()
			}
		case <-gs.ctx.Done():
			return
		}
	}
}

func (gs *GoScript) shouldRunTrigger() {
	gs.nextPeriodic = time.Now().Add(60 * time.Minute)
	for _, triggers := range gs.periodic {
		for _, t := range triggers {
			if t.GetNextTime() == nil {
				gs.Logger().Info("next time not set")
				_, err := t.NextTime(time.Now())
				if err != nil {
					gs.Logger().Error(err, "setting next time failed")
					continue
				}
			}
			if time.Now().After(*t.GetNextTime()) {
				task := gs.newTask(t, nil)
				gs.taskToRun.Add(task)

				_, err := t.NextTime(time.Now())
				if err != nil {
					gs.Logger().Error(err, "setting next time failed")
					continue
				}
			}
			if t.GetNextTime().Before(gs.nextPeriodic) {
				gs.nextPeriodic = *t.GetNextTime()
			}
		}
	}
}

func (gs *GoScript) runGronJob(gron *gronx.Gronx, start bool) {
	for expr, triggers := range gs.periodic {
		var err error
		var due bool
		if len(expr) == 0 {
			if start {
				due = true
			}
		} else {
			due, err = gron.IsDue(expr)
			if err != nil {
				gs.logger.Error(err, "gron job IsDue failed")
				continue
			}
		}
		for _, t := range triggers {
			if due {
				task := gs.newTask(t, nil)
				gs.taskToRun.Add(task)
			}
		}
	}
}

func fillNextTime(periodics map[string][]*trigger.Trigger) (time.Time, error) {
	next := time.Now().Add(60 * time.Minute)
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
