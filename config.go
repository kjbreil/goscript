package goscript

import (
	"errors"
	"github.com/goccy/go-yaml"
	mqtt "github.com/kjbreil/hass-mqtt"
	ws "github.com/kjbreil/hass-ws"
	"github.com/mitchellh/mapstructure"
	"os"
)

type Config struct {
	Websocket *ws.Config
	MQTT      *mqtt.Config
	modules   Modules
}

type Modules map[string]interface{}

var (
	ErrModuleNotFound = errors.New("module not found")
)

func (c *Config) GetModule(key string) (interface{}, error) {
	if m, ok := c.modules[key]; ok {
		return m, nil
	}
	return nil, ErrModuleNotFound
}

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
	c.modules = make(map[string]interface{})

	if err := mapstructure.Decode(configMap, &c); err != nil {
		return nil, err
	}

	for k, v := range modules {
		if m, ok := configMap[k]; ok {
			err = mapstructure.Decode(m, v)
			if err != nil {
				return nil, err
			}
			c.modules[k] = v
		}
	}

	return &c, nil
}
