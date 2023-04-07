package goscript

import (
	"errors"
	"github.com/goccy/go-yaml"
	"github.com/iancoleman/strcase"
	mqtt "github.com/kjbreil/hass-mqtt"
	ws "github.com/kjbreil/hass-ws"
	"github.com/mitchellh/mapstructure"
	"os"
	"reflect"
	"strings"
	"time"
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

	err = yaml.Unmarshal(data, &configMap)
	if err != nil {
		return nil, err
	}

	c.Modules = make(map[string]interface{})

	var decoder *mapstructure.Decoder
	decoder, err = configDecoder(&c)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(configMap)
	if err != nil {
		return nil, err
	}
	err = c.decodeModules(modules, configMap)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) decodeModules(modules Modules, configMap map[string]interface{}) error {
	for k, v := range modules {
		if m, ok := configMap[k]; ok {
			v := v
			decoder, err := configDecoder(&v)

			if err != nil {
				return err
			}
			err = decoder.Decode(m)

			if err != nil {
				return err
			}
			c.Modules[k] = v
		}
	}
	return nil
}

func configDecoder(results interface{}) (*mapstructure.Decoder, error) {
	return mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: results,
		MatchName: func(mapKey, fieldName string) bool {
			if strings.EqualFold(mapKey, fieldName) {
				return true
			}
			if strings.EqualFold(mapKey, strcase.ToSnake(fieldName)) {
				return true
			}

			return false
		},
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToTimeHookFunc(),
		),
	})
}

// stringToTimeHookFunc decodes either a simple time as am/pm or a RFC3339 formated time
func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		if strings.HasSuffix(data.(string), "pm") || strings.HasSuffix(data.(string), "am") {
			return time.Parse("15:04pm", data.(string))
		}

		if strings.HasSuffix(data.(string), "PM") || strings.HasSuffix(data.(string), "AM") {
			return time.Parse("15:04PM", data.(string))
		}

		return time.Parse(time.RFC3339, data.(string))
	}
}
