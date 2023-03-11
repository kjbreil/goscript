package goscript

import (
	"github.com/goccy/go-yaml"
	mqtt "github.com/kjbreil/hass-mqtt"
	ws "github.com/kjbreil/hass-ws"
	"github.com/mitchellh/mapstructure"
	"os"
)

type Config struct {
	Websocket *ws.Config
	MQTT      *mqtt.Config
	Modules   Modules
}

type Modules map[string]interface{}

func ParseConfig(filename string, modules Modules) (*Config, error) {
	var configMap map[string]interface{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &configMap)
	if err != nil {
		return nil, err
	}
	var c Config

	if err := mapstructure.Decode(configMap, &c); err != nil {
		return nil, err
	}

	for k, v := range modules {
		if m, ok := configMap[k]; ok {
			err = mapstructure.Decode(m, v)
			if err != nil {
				return nil, err
			}
		}
	}
	c.Modules = modules

	return &c, nil
}
