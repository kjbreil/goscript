package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kjbreil/goscript/pkg/control"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/logger"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/goscript/pkg/service"
	"github.com/kjbreil/goscript/pkg/state"
	"github.com/kjbreil/goscript/pkg/trigger"
	hassmqtt "github.com/kjbreil/hass-mqtt"
	hassws "github.com/kjbreil/hass-ws/pkg/hass"
	"github.com/kjbreil/hass-ws/model"
	"github.com/kjbreil/hass-ws/services"
)

// GoScript is the base type for GoScript holding all the state and functionality for interacting with Home Assistant
type GoScript struct {
	config *Config
	mqtt   *hassmqtt.Client
	ws     *hassws.Client
	// homekit *homekit.HomeKit

	Runner *trigger.Runner

	devices *device.Devices

	areaRegistry map[string][]model.Result

	// Context for the GoScript
	ctx    context.Context
	cancel context.CancelFunc

	ServiceChan service.Chan

	requests *control.Requests

	// states store
	states state.States

	logger *slog.Logger
}

// New creates a new GoScript instance
func New(c *Config, logger *slog.Logger) (*GoScript, error) {
	var err error

	gs := &GoScript{
		config:   c,
		logger:   logger,
		requests: control.NewRequests(),
	}
	gs.ctx, gs.cancel = context.WithCancel(context.Background())

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

	gs.ServiceChan = make(chan services.Service, 100)

	gs.Runner = trigger.NewRunner(gs.ctx, &gs.states, gs.ServiceChan, gs.logger)
	gs.devices = device.NewDevices()

	return gs, nil
}

// Connect connects to the WebSocket server and MQTT server as setup
// all options need to be passed before firing connect, anything added after will have odd effects
func (gs *GoScript) Connect() error {
	var err error

	// initialize the modules
	for _, m := range gs.config.Modules {
		err = m.Init(gs.ctx, gs.requests, gs.logger)
		if err != nil {
			return err
		}
		// for _, t := range m.Triggers() {
		// 	gs.Runner.AddTrigger(t)
		// }
		// for _, d := range m.Devices() {
		// 	err = gs.AddDevice(d)
		// }
	}

	if gs.mqtt != nil {
		err = gs.mqtt.Connect()
		if err != nil {
			if !errors.Is(err, hassmqtt.ErrNoDeviceFound) {
				return err
			}
		}
		gs.logger.Info("MQTT connected")
	}

	// for moduleName, m := range gs.config.Modules {
	// 	err = m.Run()
	// 	if err != nil {
	// 		gs.logger.Error(fmt.Sprintf("could not run module %s", moduleName), "error", err.Error())
	// 	} else {
	// 		gs.logger.Info(fmt.Sprintf("module %s running", moduleName))
	// 	}
	// }

	// start the message handler. This starts after all the modules have been initialized and run
	gs.messageHandler()

	// Add a subscription for the websocket on all events
	gs.ws.AddSubscription(model.EventTypeAll)

	// Handle all messages
	gs.ws.OnMessage = gs.handleHassMessage
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
	gs.Runner.RunPeriodic()

	// Run the modules
	for moduleName := range gs.config.Modules {
		if m, ok := gs.config.Modules[moduleName]; ok {
			go func(moduleName string, m module.Module) {
				defer func() {
					if r := recover(); r != nil {
						gs.logger.Error(fmt.Sprintf("module %s panicked", moduleName), "panic", r)
					}
				}()
				err := m.Run()
				if err != nil {
					gs.logger.Error(fmt.Sprintf("could not run module %s", moduleName), "error", err.Error())
				} else {
					gs.logger.Info(fmt.Sprintf("module %s running", moduleName))
				}
			}(moduleName, m)
		}

	}

	// homekit integration needs to be setup after all modules have been run because devices cannot be added to homekit
	// after starting

	// if gs.config.Homekit != nil {
	// 	gs.homekit = homekit.New(gs.ctx)
	// }

	gs.logger.Info("GoScript started")

	return nil
}

func (gs *GoScript) CallService(s services.Service) *hassws.Response {
	return gs.ws.CallService(s)
}

// Logger returns the logr to create your own logs
func (gs *GoScript) Logger() *slog.Logger {
	return gs.logger
}

func (gs *GoScript) runFunctions() {
	defer func() {
		gs.logger.Info("runFunctions exited")
	}()
	timer := time.NewTicker(10 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-gs.ctx.Done():
			return
		case <-timer.C:
			for _, t := range gs.Runner.TaskToRun() {
				go gs.Runner.RunTask(t)
			}
		}
	}
}

// Close the connections to WebSocket and MQTT
func (gs *GoScript) Close() {
	for _, m := range gs.config.Modules {
		m.Close()
	}
	gs.cancel()
	err := gs.ws.Close()
	gs.mqtt.Disconnect()
	if err != nil {
		gs.logger.Error("error closing websocket", "error", err.Error())
	}
}

// GetModule returns the config module in interface{} form, must be cast to module type
func (gs *GoScript) GetModule(key string) (interface{}, error) {
	return gs.config.GetModule(key)
}

func GetModule[T any](gs *GoScript, key string) T {
	if v, ok := gs.config.Modules[key]; ok {
		return v.(T)
	}
	panic(ErrModuleNotFound)
}

func DefaultLogger() *slog.Logger {
	return slog.New(logger.NewHandler(os.Stdout, nil))
}

// DefaultLoggerWithLevel returns a logger with the specified level
func DefaultLoggerWithLevel(level slog.Level) *slog.Logger {
	return slog.New(logger.NewHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
