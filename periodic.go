package goscript

import (
	"context"
	"github.com/adhocore/gronx"
	"github.com/google/uuid"
	"log"
	"time"
)

type Periodic []string

func Periodics(times ...string) []string {
	rtn := make([]string, len(times))
	//TODO: Validate the cron strings
	for i := range times {
		rtn[i] = times[i]
	}
	return rtn
}

func (gs *GoScript) runPeriodic() {
	gron := gronx.New()
	// try and run the jobs right away this will run the every minute ones and at start ones
	gs.runGronJob(&gron, true)
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
			gs.runGronJob(&gron, false)
		case <-gs.ctx.Done():
			return
		}
	}

}

func (gs *GoScript) runGronJob(gron *gronx.Gronx, start bool) {
	funcToRun := make(map[uuid.UUID]*Task)

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
				// TODO: Need error bus
				log.Println(err)
				continue
			}
		}
		for _, t := range triggers {

			if due {
				task := &Task{
					Message:     nil,
					States:      gs.GetStates(t.States),
					gs:          gs,
					states:      t.States,
					f:           t.Func,
					waitRequest: make(chan *Trigger),
					waitDone:    make(chan bool),
				}
				if t.Unique != nil {
					t.Unique.cancel()
					t.Unique.ctx, t.Unique.cancel = context.WithCancel(context.Background())
					task.ctx, task.cancel = t.Unique.ctx, t.Unique.cancel
					funcToRun[t.uuid] = task
				} else {
					task.ctx, task.cancel = context.WithCancel(context.Background())
					funcToRun[uuid.New()] = task
				}
			}
		}
	}

	for _, t := range funcToRun {
		go t.f(t)
		go gs.taskWaitRequest(t)
	}
}
