package goscript

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/google/uuid"
	hassmqtt "github.com/kjbreil/hass-mqtt"
	hassws "github.com/kjbreil/hass-ws"
	"github.com/kjbreil/hass-ws/model"
	"github.com/kjbreil/hass-ws/services"
	"sync"
	"time"
)

// GoScript is the base type for GoScript holding all the state and functionality for interacting with Home Assistant
type GoScript struct {
	config *Config
	mqtt   *hassmqtt.Client
	ws     *hassws.Client

	// maps holding state based triggers
	periodic      map[string][]*Trigger
	nextPeriodic  time.Time
	triggers      map[string][]*Trigger
	domainTrigger map[string][]*Trigger

	devices map[string]*Device

	areaRegistry map[string][]model.Result

	taskToRun taskMap

	// Context for the GoScript
	ctx    context.Context
	cancel context.CancelFunc

	ServiceChan ServiceChan
	// states store
	states States

	logger logr.Logger
}

// ServiceChan is a channel to send services to be run to
type ServiceChan chan services.Service

// New creates a new GoScript instance
func New(c *Config, logger logr.Logger) (*GoScript, error) {
	var err error

	gs := &GoScript{
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

	gs.triggers = make(map[string][]*Trigger)
	gs.domainTrigger = make(map[string][]*Trigger)
	gs.periodic = make(map[string][]*Trigger)
	gs.ServiceChan = make(chan services.Service, 100)
	gs.taskToRun = taskMap{
		tasks: make(map[uuid.UUID]*Task),
		m:     &sync.Mutex{},
	}

	gs.states = States{
		s: make(map[string]*State),
		m: &sync.Mutex{},
	}
	gs.devices = make(map[string]*Device)

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

	go gs.runService()

	go gs.runPeriodic()

	gs.logger.Info("GoScript started")

	return nil
}

// Logger returns the logr to create your own logs
func (gs *GoScript) Logger() logr.Logger {
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

func (gs *GoScript) runFunctions() {
	defer func() {
		gs.logger.Info("runFunctions exited")
	}()
	timer := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-gs.ctx.Done():
			return
		case <-timer.C:
			var ran []uuid.UUID
			gs.taskToRun.m.Lock()
			for u, t := range gs.taskToRun.tasks {
				go gs.runTask(t)
				ran = append(ran, u)
			}
			for _, u := range ran {
				delete(gs.taskToRun.tasks, u)
			}
			gs.taskToRun.m.Unlock()
		}
	}
}

// Close the connections to WebSocket and MQTT
func (gs *GoScript) Close() {
	gs.cancel()
	err := gs.ws.Close()
	gs.mqtt.Disconnect()
	if err != nil {
		gs.logger.Error(err, "error closing websocket")
	}
}

// GetModule returns the config module in interface{} form, must be cast to module type
func (gs *GoScript) GetModule(key string) (interface{}, error) {
	return gs.config.GetModule(key)
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
