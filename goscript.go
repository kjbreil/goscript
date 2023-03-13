package goscript

import (
	"context"
	"errors"
	"github.com/go-logr/logr"
	hass_mqtt "github.com/kjbreil/hass-mqtt"
	hass_ws "github.com/kjbreil/hass-ws"
	"github.com/kjbreil/hass-ws/model"
	"github.com/kjbreil/hass-ws/services"
	"time"
)

type GoScript struct {
	config *Config
	mqtt   *hass_mqtt.Client
	ws     *hass_ws.Client

	periodic map[string][]*Trigger
	triggers map[string][]*Trigger

	ctx    context.Context
	cancel context.CancelFunc

	ServiceChan ServiceChan
	states      states

	logger logr.Logger
}
type ServiceChan chan services.Service

// New creates a new GoScript instance
func New(c *Config) (*GoScript, error) {
	var err error

	gs := &GoScript{
		config: c,
	}

	gs.mqtt, err = hass_mqtt.NewClient(*gs.config.MQTT)
	if err != nil {
		return nil, err
	}

	gs.ws, err = hass_ws.NewClient(gs.config.Websocket)
	if err != nil {
		return nil, err
	}

	gs.triggers = make(map[string][]*Trigger)
	gs.periodic = make(map[string][]*Trigger)
	gs.ServiceChan = make(chan services.Service, 100)

	return gs, nil
}

// Connect connects to the WebSocket server and MQTT server as setup
// all options need to be passed before firing connect, anything added after will have odd effects
func (gs *GoScript) Connect() error {
	var err error

	gs.ctx, gs.cancel = context.WithCancel(context.Background())

	if gs.mqtt != nil {
		err = gs.mqtt.Connect()
		if err != nil {
			if !errors.Is(err, hass_mqtt.ErrNoDeviceFound) {
				return err
			}
		}
		gs.logger.Info("MQTT connected")
	}

	// Add a subscription for the websocket on all events
	gs.ws.AddSubscription(model.EventTypeAll)

	// Handle all messages
	gs.ws.OnMessage = gs.handleMessage
	// handle running get states
	gs.ws.OnGetState = gs.handleGetStates
	// setup hass_ws to initialize all states at connect. This is run through the triggers.
	gs.ws.InitStates = true

	err = gs.ws.Connect()
	if err != nil {
		return err
	}
	gs.logger.Info("Websocket connected")

	time.Sleep(100 * time.Millisecond)

	go gs.runService()

	go gs.runPeriodic()

	gs.logger.Info("GoScript started")

	return nil
}

func (gs *GoScript) Logger(logger logr.Logger) {
	gs.logger = logger
}
func (gs *GoScript) GetLogger() logr.Logger {
	return gs.logger
}

func (gs *GoScript) runService() {
	for {
		select {
		case <-gs.ctx.Done():
			return
		case s := <-gs.ServiceChan:
			gs.ws.CallService(s)
		}
	}
}

func (gs *GoScript) Close() {
	gs.cancel()
	err := gs.ws.Close()
	if err != nil {
		gs.logger.Error(err, "error closing websocket")
	}
}

func (gs *GoScript) GetModule(key string) (interface{}, error) {
	return gs.config.GetModule(key)
}
