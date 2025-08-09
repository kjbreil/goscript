package core

import (
	"fmt"
	"time"

	"github.com/kjbreil/goscript/pkg/control"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/hass-mqtt/common"
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

func (gs *GoScript) handleMessage(m control.Request) error {
	// create the response object to use in the callback
	rsp := control.Response{
		To:   m.From,
		From: "goscript",
	}

	if m.Device != nil {
		err := gs.AddDevice(m.Device)
		if err != nil {
			return err
		}
		// go func() {
		// 	err = gs.homekit.Run(gs.devices.GetHomekitAccessories())
		// 	if err != nil {
		// 		gs.logger.Error(err.Error())
		// 	}
		// }()
	}

	if m.Trigger != nil {
		gs.Runner.AddTrigger(m.Trigger)
	}

	if m.TaskTrigger != nil && m.MQTTMessage != nil {
		task := gs.Runner.NewTask(m.TaskTrigger, nil)
		task.MqttMessage = *m.MQTTMessage
		gs.Runner.AddTask(task)
	}

	// send a service call
	if m.Service != nil {
		serviceRespAwait := gs.CallService(*m.Service)
		serviceResp := serviceRespAwait.Timeout(time.Second * 5)
		rsp.ServiceRsp = serviceResp
	}

	if m.GetStates != nil {
		states := gs.states.SubSet(m.GetStates)
		rsp.States = &states
	}

	if m.History != nil {
		histories, err := gs.GetHistory(m.History.Start, m.History.End, m.History.Entities...)
		if err != nil {
			return err
		}
		rsp.Histories = &histories
	}

	if m.MQTTPublish != nil {
		token := gs.mqtt.Publish(m.MQTTPublish.Topic, m.MQTTPublish.QoS, m.MQTTPublish.Retained, m.MQTTPublish.Payload)

		token.WaitTimeout(common.WaitTimeout)

		if token.Error() != nil {
			return token.Error()
		}
	}

	if m.Log != nil {
		msg := fmt.Sprintf("(%s) %s", m.From, m.Log.Msg)

		if m.Log.Err != nil {
			gs.Logger().Error(msg, "error", m.Log.Err, "caller", m.Log.Caller)
		} else {
			gs.Logger().Info(msg, "caller", m.Log.Caller)
		}
	}

	if m.Module != nil {
		switch mod := m.Module.(type) {
		case module.Module:
			if c, ok := gs.config.Modules.Get(mod.Name()); ok {
				rsp.Module = c
			}
		}
	}

	if m.Callback != nil {
		go func() {
			err := m.Callback(rsp)
			if err != nil {
				gs.logger.Error("callback returned an error", "err", err.Error())
			}
		}()
	}

	return nil
}

func (gs *GoScript) sendResponse(rsp control.Response) {
	if mod, ok := gs.config.Modules.Get(rsp.To); ok {
		mod.Responses(rsp)
	}
}
