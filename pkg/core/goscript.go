package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/service"
	"github.com/kjbreil/goscript/pkg/state"
	"github.com/kjbreil/goscript/pkg/trigger"
	hassmqtt "github.com/kjbreil/hass-mqtt"
	hassws "github.com/kjbreil/hass-ws"
	"github.com/kjbreil/hass-ws/model"
	"github.com/kjbreil/hass-ws/services"
	"time"
)

// Core is the base type for Core holding all the state and functionality for interacting with Home Assistant
type Core struct {
	config *Config
	mqtt   *hassmqtt.Client
	ws     *hassws.Client

	TrigRunner *trigger.Runner

	devices map[string]*device.Device

	areaRegistry map[string][]model.Result

	// Context for the Core
	ctx    context.Context
	cancel context.CancelFunc

	ServiceChan service.Chan
	// states store
	states state.States

	logger logr.Logger
}

// New creates a new Core instance
func New(c *Config, logger logr.Logger) (*Core, error) {
	var err error

	gs := &Core{
		config: c,
		logger: logger,
	}

	gs.mqtt, err = hassmqtt.NewClientWithLogger(*gs.config.MQTT, gs.logger)
	if err != nil {
		return nil, err
	}

	gs.ws, err = hassws.NewClientWithLogger(gs.config.Websocket, gs.logger)
	if err != nil {
		return nil, err
	}
	gs.ws.Logger()
	gs.states = state.NewStates()

	gs.TrigRunner = trigger.NewRunner(gs.ctx, &gs.states, gs.ServiceChan, gs.logger)

	gs.ServiceChan = make(chan services.Service, 100)

	gs.devices = make(map[string]*device.Device)

	return gs, nil
}

// Connect connects to the WebSocket server and MQTT server as setup
// all options need to be passed before firing connect, anything added after will have odd effects
func (gs *Core) Connect() error {
	var err error

	gs.ctx, gs.cancel = context.WithCancel(context.Background())

	if gs.mqtt != nil {
		err = gs.mqtt.Connect()
		if err != nil {
			if !errors.Is(err, hassmqtt.ErrNoDeviceFound) {
				return err
			}
		}
		gs.logger.Info("MQTT connected")
	}

	// Add a subscription for the websocket on all events
	gs.ws.AddSubscription(model.EventTypeAll)

	// Handle all messages
	gs.ws.OnMessage = gs.handleMessage
	// handle running get States
	gs.ws.OnGetState = gs.handleGetStates
	// setup hass_ws to initialize all States at connect. This is run through the triggers.
	gs.ws.InitStates = true

	err = gs.ws.Connect()
	if err != nil {
		return err
	}
	gs.logger.Info("Websocket connected")

	gs.fillAreaRegistry()

	time.Sleep(100 * time.Millisecond)

	go gs.runFunctions()

	service.Run(gs.ctx, gs.ws, gs.logger, gs.ServiceChan)

	// TODO: Change this into a RUN function passing the periodics
	gs.TrigRunner.RunPeriodic()

	gs.logger.Info("Core started")

	return nil
}

// Logger returns the logr to create your own logs
func (gs *Core) Logger() logr.Logger {
	return gs.logger
}

func (gs *Core) runFunctions() {
	defer func() {
		gs.logger.Info("runFunctions exited")
	}()
	timer := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-gs.ctx.Done():
			return
		case <-timer.C:
			for _, t := range gs.TrigRunner.TaskToRun() {
				go gs.TrigRunner.RunTask(t)
			}
		}
	}
}

// Close the connections to WebSocket and MQTT
func (gs *Core) Close() {
	gs.cancel()
	err := gs.ws.Close()
	gs.mqtt.Disconnect()
	if err != nil {
		gs.logger.Error(err, "error closing websocket")
	}
}

// GetModule returns the config module in interface{} form, must be cast to module type
func (gs *Core) GetModule(key string) (interface{}, error) {
	return gs.config.GetModule(key)
}

func GetModule[T any](gs *Core, key string) T {
	if v, ok := gs.config.Modules[key]; ok {
		return v.(T)
	}
	panic(ErrModuleNotFound)
}

func DefaultLogger() logr.Logger {
	log := funcr.New(
		func(pfx, args string) { fmt.Println(pfx, args) },
		funcr.Options{
			LogCaller:    funcr.None,
			LogTimestamp: true,
			Verbosity:    1,
		})
	return log.WithName("goscript")
}
