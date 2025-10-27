package module

import (
	"runtime"
	"strconv"

	"github.com/kjbreil/goscript/pkg/control"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
)

func SendDevices(m Module, ds *device.Devices) {
	for _, d := range ds.Slice() {
		m.Requests().Chan() <- control.Request{
			To:      "goscript",
			From:    m.Name(),
			Trigger: nil,
			Device:  d,
		}
	}
}

func SendDevice(m Module, d *device.Device) {
	m.Requests().Chan() <- control.Request{
		To:      "goscript",
		From:    m.Name(),
		Trigger: nil,
		Device:  d,
	}
}

func SendTriggers(m Module, triggers trigger.Triggers) {
	for _, t := range triggers {
		m.Requests().Chan() <- control.Request{
			To:      "goscript",
			From:    m.Name(),
			Trigger: t,
		}
	}
}

func SendPublish(m Module, mqttPublish *control.MQTTPublish) {
	m.Requests().Chan() <- control.Request{
		To:          "goscript",
		From:        m.Name(),
		MQTTPublish: mqttPublish,
	}
}

func SendInfo(m Module, msg string) {
	var caller string
	_, file, line, ok := runtime.Caller(1)
	if ok {
		caller = file + ":" + strconv.Itoa(line)
	}
	m.Requests().Chan() <- control.Request{
		To:   "goscript",
		From: m.Name(),
		Log: &control.Log{
			Msg:    msg,
			Err:    nil,
			Caller: caller,
		},
	}
}
func SendErr(m Module, err error, msg string) {
	var caller string
	_, file, line, ok := runtime.Caller(1)
	if ok {
		caller = file + ":" + strconv.Itoa(line)
	}
	m.Requests().Chan() <- control.Request{
		To:   "goscript",
		From: m.Name(),
		Log: &control.Log{
			Msg:    msg,
			Err:    err,
			Caller: caller,
		},
	}
}

func SendService(m Module, s services.Service) {
	m.Requests().Chan() <- control.Request{
		To:      "goscript",
		From:    m.Name(),
		Service: &s,
	}
}

func SendServiceCallback(m Module, s services.Service, callback func(rsp control.Response) error) {
	m.Requests().Chan() <- control.Request{
		To:       "goscript",
		From:     m.Name(),
		Service:  &s,
		Callback: callback,
	}
}
