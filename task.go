package goscript

import (
	"context"
	"time"
)

type Task struct {
	ctx    context.Context
	cancel context.CancelFunc
	states []string
	f      TriggerFunc
}

func (t *Task) Sleep(dur time.Duration) bool {
	timer := time.NewTimer(dur)
	select {
	case <-timer.C:
		return true
	case <-t.ctx.Done():
		return false
	}
}
