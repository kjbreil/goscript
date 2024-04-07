package core

import (
	"github.com/kjbreil/goscript/pkg/module"
)

func (gs *GoScript) messageHandler() {
	go func() {
		for {
			select {
			case <-gs.ctx.Done():
				return
			case m := <-gs.requests.Chan():
				if m.To == "goscript" {
					err := gs.handleMessage(m)
					if err != nil {
						gs.Logger().Error(err.Error())
					}
				} else {
					panic("not done message handler")
				}
			}
		}
	}()

}

func (gs *GoScript) handleMessage(m module.Request) error {

	if m.Device != nil {
		err := gs.AddDevice(m.Device)
		if err != nil {
			return err
		}
	}

	if m.Trigger != nil {
		gs.Runner.AddTrigger(m.Trigger)
	}

	if m.TaskTrigger != nil && m.Message != nil {
		task := gs.Runner.NewTask(m.TaskTrigger, nil)
		task.MqttMessage = *m.Message
		gs.Runner.AddTask(task)
	}

	if m.Service != nil {
		gs.ServiceChan <- *m.Service
	}

	if m.Log != nil {
		if m.Log.Err != nil {
			gs.Logger().Error(m.Log.Msg, "error", m.Log.Err, "caller", m.Log.Caller)
		} else {
			gs.Logger().Info(m.Log.Msg, "caller", m.Log.Caller)
		}
	}

	if m.Module != "" {
		if c, ok := gs.config.Modules.Get(m.Module); ok {

			gs.sendResponse(module.Response{
				To:     m.From,
				From:   "goscript",
				Module: c,
			})
		}
	}

	return nil
}

func (gs *GoScript) sendResponse(rsp module.Response) {
	if mod, ok := gs.config.Modules.Get(rsp.To); ok {
		mod.Responses(rsp)
	}
}
