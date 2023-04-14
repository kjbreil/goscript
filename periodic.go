package goscript

import (
	"github.com/adhocore/gronx"
	"time"
)

// Periodic is the list of cron expressions to run periodically
type Periodic []string

// Periodics is a helper function to add multiple strings without needing a []string{}
func Periodics(times ...string) []string {
	rtn := make([]string, len(times))
	//TODO: Validate the cron strings
	copy(rtn, times)
	return rtn
}

func (gs *GoScript) runPeriodic() {
	// TODO: Validate Periodic slice
	// run zero length immediate periodics and delete from periodic list
	for _, triggers := range gs.periodic {
		for _, t := range triggers {
			for i := range t.Periodic {
				if len(t.Periodic[i]) == 0 {
					task := gs.newTask(t, nil)
					gs.funcToRun[task.uuid] = task

					t.Periodic = append(t.Periodic[:i], t.Periodic[i+1:]...)
				}
			}
		}
	}
	delete(gs.periodic, "")

	// setup the next fire time for all triggers
	gs.fillNextTime()

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
			if t.nextTime == nil {
				gs.Logger().Info("next time not set")
				_, err := t.NextTime(time.Now())
				if err != nil {
					gs.Logger().Error(err, "setting next time failed")
					continue
				}
			}
			if time.Now().After(*t.nextTime) {
				task := gs.newTask(t, nil)
				gs.funcToRun[task.uuid] = task
				_, err := t.NextTime(time.Now())
				if err != nil {
					gs.Logger().Error(err, "setting next time failed")
					continue
				}
			}
			if t.nextTime.Before(gs.nextPeriodic) {
				gs.nextPeriodic = *t.nextTime
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
				gs.funcToRun[task.uuid] = task
			}
		}
	}
}

func (gs *GoScript) fillNextTime() {
	gs.nextPeriodic = time.Now().Add(60 * time.Minute)
	for _, triggers := range gs.periodic {
		for _, t := range triggers {
			nt, err := t.NextTime(time.Now())
			if err != nil {
				gs.Logger().Error(err, "setting next time failed")
			}
			if nt != nil && nt.Before(gs.nextPeriodic) {
				gs.nextPeriodic = *nt
			}
		}
	}
}
