// Package control provides request/response messaging for the control system.
package control

import (
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/history"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
)

const (
	requestChannelSize    = 1000
	backpressureThreshold = 800 // 80% of channel capacity
)

// Requests manages request channels and logging for the control system.
type Requests struct {
	channel chan Request
	logger  *slog.Logger
}

// NewRequests creates a new Requests instance without a logger.
func NewRequests() *Requests {
	//nolint:exhaustruct // logger is optional and may be nil
	return &Requests{
		channel: make(chan Request, requestChannelSize),
	}
}

// NewRequestsWithLogger creates a new Requests instance with a logger.
func NewRequestsWithLogger(logger *slog.Logger) *Requests {
	return &Requests{
		channel: make(chan Request, requestChannelSize),
		logger:  logger,
	}
}

// MQTTPublish represents an MQTT publish request.
type MQTTPublish struct {
	Topic    string
	QoS      byte
	Retained bool
	Payload  any
}

// Request represents a control system request with routing and callback information.
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

// Log represents logging information for a request.
type Log struct {
	Msg    string
	Err    error
	Caller string
}

// Chan returns the request channel.
func (r *Requests) Chan() chan Request {
	return r.channel
}

// Send sends a request with backpressure warning.
func (r *Requests) Send(req Request) {
	// Warn if channel is getting full (>80% capacity)
	if r.logger != nil && len(r.channel) > backpressureThreshold {
		r.logger.Warn("request channel backpressure detected", "buffered", len(r.channel), "capacity", cap(r.channel))
	}
	r.channel <- req
}
