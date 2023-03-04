package goscript

import (
	"context"
	hass_mqtt "github.com/kjbreil/hass-mqtt"
	hass_ws "github.com/kjbreil/hass-ws"
	"github.com/kjbreil/hass-ws/model"
	"github.com/kjbreil/hass-ws/services"
	"time"
)

type GoScript struct {
	config   *Config
	mqtt     *hass_mqtt.Client
	ws       *hass_ws.Client
	triggers map[string][]*Trigger

	ctx    context.Context
	cancel context.CancelFunc

	ServiceChan chan services.Service
	states      states
}

func New(c *Config) (*GoScript, error) {
	var err error

	gs := &GoScript{
		config: c,
	}

	gs.mqtt, err = hass_mqtt.NewClient(gs.config.MQTT)
	if err != nil {
		return nil, err
	}

	gs.ws, err = hass_ws.NewClient(&gs.config.Websocket)
	if err != nil {
		return nil, err
	}

	gs.triggers = make(map[string][]*Trigger)
	gs.ServiceChan = make(chan services.Service, 100)

	return gs, nil
}

func (gs *GoScript) Connect() error {
	var err error

	gs.ctx, gs.cancel = context.WithCancel(context.Background())

	//err = gs.mqtt.Connect()
	//if err != nil {
	//	return err
	//}

	gs.ws.AddSubscription(model.EventTypeAll)

	gs.ws.OnMessage = gs.handleMessage

	gs.ws.OnGetState = gs.handleGetStates

	err = gs.ws.Connect()
	if err != nil {
		return err
	}

	gs.ws.GetStates()

	time.Sleep(1 * time.Second)

	go gs.runService()

	return nil
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
	gs.ws.Close()
}