package goscript

import (
	"github.com/goccy/go-yaml"
	mqtt "github.com/kjbreil/hass-mqtt"
	ws "github.com/kjbreil/hass-ws"
	"os"
)

type Config struct {
	Websocket ws.Config
	MQTT      mqtt.Config
}

func ParseConfig(filename string) (*Config, error) {

	var c Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
