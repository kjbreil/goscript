package goscript

import (
	"errors"
	"github.com/goccy/go-yaml"
	"github.com/iancoleman/strcase"
	mqtt "github.com/kjbreil/hass-mqtt"
	ws "github.com/kjbreil/hass-ws"
	"github.com/mitchellh/mapstructure"
	"os"
	"strings"
)

type Config struct {
	Websocket *ws.Config
	MQTT      *mqtt.Config
	Modules   Modules
}

type Modules map[string]interface{}

var (
	ErrModuleNotFound = errors.New("module not found")
)

func (c *Config) GetModule(key string) (interface{}, error) {
	if m, ok := c.Modules[key]; ok {
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
	var c Config
	//err = yaml.Unmarshal(data, &c)
	//if err != nil {
	//	return nil, err
	//}

	err = yaml.Unmarshal(data, &configMap)
	if err != nil {
		return nil, err
	}
	//var c Config
	c.Modules = make(map[string]interface{})

	if err := mapstructure.Decode(configMap, &c); err != nil {
		return nil, err
	}

	for k, v := range modules {
		if m, ok := configMap[k]; ok {
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Result: v,
				MatchName: func(mapKey, fieldName string) bool {
					if strings.EqualFold(mapKey, fieldName) {
						return true
					}
					if strings.EqualFold(mapKey, strcase.ToSnake(fieldName)) {
						return true
					}

					return false
				},
			})

			if err != nil {
				return nil, err
			}
			err = decoder.Decode(m)

			//err = mapstructure.Decode(m, v)
			if err != nil {
				return nil, err
			}
			c.Modules[k] = v
		}
	}

	return &c, nil
}
