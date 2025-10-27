package control

import (
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/history"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
)

type Requests struct {
	channel chan Request
	logger  *slog.Logger
}

func NewRequests() *Requests {
	return &Requests{
		channel: make(chan Request, 1000),
	}
}

func NewRequestsWithLogger(logger *slog.Logger) *Requests {
	return &Requests{
		channel: make(chan Request, 1000),
		logger:  logger,
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
	Device      *device.Device
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

// Send sends a request with backpressure warning
func (r *Requests) Send(req Request) {
	// Warn if channel is getting full (>80% capacity)
	if r.logger != nil && len(r.channel) > 800 {
		r.logger.Warn("request channel backpressure detected", "buffered", len(r.channel), "capacity", cap(r.channel))
	}
	r.channel <- req
}
