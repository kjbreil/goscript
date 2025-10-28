package trigger

import (
	"context"
	"testing"
	"time"
)

func TestTask_Sleep(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		ctx:    ctx,
		cancel: cancel,
	}

	go func() {
		task.Sleep(100 * time.Millisecond)
		task.cancel()
	}()

	select {
	case <-ctx.Done():
		return
	case <-time.After(101 * time.Millisecond):
		t.Error("context should have been cancelled")
	}
}

func TestTask_SleepCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		ctx:    ctx,
		cancel: cancel,
	}

	done := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("task exited")
				done <- struct{}{}
			}
		}()
		task.Sleep(100 * time.Millisecond)
		t.Error("Should not have reached here")
	}()

	// cancel the task
	task.cancel()
	<-done
}
