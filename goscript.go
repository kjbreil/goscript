package goscript

import (
	"github.com/hashicorp/go-memdb"
	hass_mqtt "github.com/kjbreil/hass-mqtt"
	hass_ws "github.com/kjbreil/hass-ws"
	"github.com/kjbreil/hass-ws/model"
	"time"
)

type GoScript struct {
	config   *Config
	db       *memdb.MemDB
	mqtt     *hass_mqtt.Client
	ws       *hass_ws.Client
	triggers map[string][]*Trigger

	states map[string]*State
}

func New(c *Config) (*GoScript, error) {
	var err error

	gs := &GoScript{
		config: c,
	}

	gs.db, err = memdb.NewMemDB(memdbSchema())
	if err != nil {
		return nil, err
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

	return gs, nil
}

func (gs *GoScript) Connect() error {
	var err error

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

	return nil
}

func (gs *GoScript) Close() {
	gs.ws.Close()
}
