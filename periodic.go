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
	gron := gronx.New()
	// try and run the jobs right away this will run the every minute ones and at start ones
	go gs.runGronJob(&gron, true)
	// Wait until the next whole minute to start the ticker
	now := time.Now()
	nowMinute, _ := time.Parse("2006-01-02T15:04Z07:00", now.Format("2006-01-02T15:04Z07:00"))
	nextMinute := nowMinute.Add(time.Minute)
	dur := nextMinute.Sub(now)
	time.Sleep(dur)

	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			go gs.runGronJob(&gron, false)
		case <-gs.ctx.Done():
			return
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
