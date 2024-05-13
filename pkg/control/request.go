package control

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/history"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
)

type Requests struct {
	channel chan Request
}

func NewRequests() *Requests {
	return &Requests{
		channel: make(chan Request, 1000),
	}
}

type MQTTPublish struct {
	Topic    string
	QoS      byte
	Retained bool
	Payload  any
}

type Request struct {
	To   string
	From string

	Trigger     *trigger.Trigger
	TaskTrigger *trigger.Trigger
	MQTTMessage *mqtt.Message
	MQTTPublish *MQTTPublish
	Device      *device.GSDevice
	Service     *services.Service
	History     *history.GetHistories
	GetStates   []string

	Callback func(rsp Response) error

	Log *Log

	Module any
}

type Log struct {
	Msg    string
	Err    error
	Caller string
}

func (r *Requests) Chan() chan Request {
	return r.channel
}
