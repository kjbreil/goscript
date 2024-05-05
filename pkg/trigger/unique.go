package trigger

import (
	"context"
	"github.com/google/uuid"
)

// The Unique task will wait until the currently running task finishes to start. To quickly kill tasks that are
// running it is important to only use the task methods within a trigger function instead of time.Sleep. If
// Unique.KillMe is set to true the task will not be setup and will not run if another task is running of the same type.
// Unique.UUID is used to link multiple triggers together. For example two triggers that control the same light and you
// only want one of the trigger functions to run at a time. Unique.Wait waits until the current task is finished before
// running, will build up multiple tasks. The queue is based on the UUID so linking the UUID's will make one bit queue.
type Unique struct {
	KillMe bool
	Wait   bool
	UUID   *uuid.UUID

	running *bool
	ctx     context.Context
	cancel  context.CancelFunc
}

func (u *Unique) Running() *bool {
	if u.running == nil {
		return ptr(false)
	}
	return u.running
}

func (u *Unique) SetRunning(running bool) {
	u.running = &running
}

func (u *Unique) NewCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	u.ctx, u.cancel = context.WithCancel(ctx)
	return u.ctx, u.cancel
}
func (u *Unique) CancelNew(ctx context.Context) (context.Context, context.CancelFunc) {
	if u.cancel != nil {
		u.cancel()
	}

	u.ctx, u.cancel = context.WithCancel(ctx)
	return u.ctx, u.cancel
}
