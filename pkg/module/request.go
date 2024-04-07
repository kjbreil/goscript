package module

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
	"runtime"
	"strconv"
)

type Requests struct {
	channel chan Request
}

type Request struct {
	To   string
	From string

	Trigger     *trigger.Trigger
	TaskTrigger *trigger.Trigger
	Message     *mqtt.Message
	Device      *device.Device
	Service     *services.Service

	Log    *Log
	Module string
}

type Log struct {
	Msg    string
	Err    error
	Caller string
}

func (r *Requests) Chan() chan Request {
	return r.channel
}

func NewRequests() *Requests {
	return &Requests{
		channel: make(chan Request, 1000),
	}
}

func SendDevices(m Module, ds device.Devices) {
	for _, d := range ds {
		m.Requests().channel <- Request{
			To:      "goscript",
			From:    m.Name(),
			Trigger: nil,
			Device:  d,
		}
	}
}

func SendTriggers(m Module, triggers trigger.Triggers) {
	for _, t := range triggers {
		m.Requests().channel <- Request{
			To:      "goscript",
			From:    m.Name(),
			Trigger: t,
		}
	}
}

func SendInfo(m Module, msg string) {
	var caller string
	_, file, line, ok := runtime.Caller(1)
	if ok {
		caller = file + ":" + strconv.Itoa(line)
	}
	m.Requests().channel <- Request{
		To:   "goscript",
		From: m.Name(),
		Log: &Log{
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
	m.Requests().channel <- Request{
		To:   "goscript",
		From: m.Name(),
		Log: &Log{
			Msg:    msg,
			Err:    err,
			Caller: caller,
		},
	}
}

func SendService(m Module, s services.Service) {
	m.Requests().channel <- Request{
		To:      "goscript",
		From:    m.Name(),
		Service: &s,
	}
}

func TaskMQTT(m Module, tr *trigger.Trigger) func(message mqtt.Message, client mqtt.Client) {
	// setup the trigger
	tr = trigger.SetupTrigger(tr)

	return func(message mqtt.Message, client mqtt.Client) {
		m.Requests().channel <- Request{
			To:          "goscript",
			From:        m.Name(),
			TaskTrigger: tr,
			Message:     &message,
		}

	}
}
