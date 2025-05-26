package module

import (
	"context"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kjbreil/goscript/pkg/control"
	"github.com/kjbreil/goscript/pkg/trigger"
)

type Module interface {
	// Init brings in data streams that might be needed and returns the triggers the module provides
	Init(ctx context.Context, requests *control.Requests, logger *slog.Logger) error
	Responses(rsp control.Response)
	Requests() *control.Requests
	Run() error
	Update() error
	Name() string
	Close() error

	mustImplementBase()
}

type Base struct {
	Ctx      context.Context
	requests *control.Requests
	logger   *slog.Logger
}

func (m *Base) mustImplementBase() {}

func (m *Base) Init(ctx context.Context, requests *control.Requests, logger *slog.Logger) error {
	m.AssignBase(ctx, requests)
	m.logger = logger
	return nil
}

func (m *Base) AssignBase(ctx context.Context, requests *control.Requests) {
	m.Ctx = ctx
	m.requests = requests
}

func (m *Base) Responses(_ control.Response) {

}
func (m *Base) Requests() *control.Requests {
	return m.requests
}

func TaskMQTT(m Module, tr *trigger.Trigger) func(message mqtt.Message, client mqtt.Client) {
	// setup the trigger
	tr = trigger.SetupTrigger(tr)

	return func(message mqtt.Message, _ mqtt.Client) {
		m.Requests().Chan() <- control.Request{
			To:          "goscript",
			From:        m.Name(),
			TaskTrigger: tr,
			MQTTMessage: &message,
		}
	}
}
