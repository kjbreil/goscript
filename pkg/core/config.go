package core

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/iancoleman/strcase"
	"github.com/kjbreil/goscript/pkg/homekit"
	"github.com/kjbreil/goscript/pkg/module"
	mqtt "github.com/kjbreil/hass-mqtt"
	ws "github.com/kjbreil/hass-ws/pkg/hass"
	"github.com/mitchellh/mapstructure"
)

type GoScriptConfig struct {
	LogLevel string
}

// GetLogLevel parses the LogLevel string and returns the corresponding slog.Level.
// Supported values: "debug", "info", "warn", "error" (case-insensitive).
// Returns slog.LevelInfo if the value is empty or invalid.
func (g *GoScriptConfig) GetLogLevel() slog.Level {
	if g == nil || g.LogLevel == "" {
		return slog.LevelInfo
	}

	switch strings.ToLower(g.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type Config struct {
	GoScript  *GoScriptConfig
	Websocket *ws.Config
	MQTT      *mqtt.Config
	Homekit   *homekit.HomeKitConfig
	Modules   module.Modules
	Timezone  string
}

var (
	ErrModuleNotFound = errors.New("module not found")
)

func (c *Config) GetModule(key string) (interface{}, error) {
	if m, ok := c.Modules[key]; ok {
		return m, nil
	}
	return nil, ErrModuleNotFound
}

func ParseConfig(filename string, modules []module.Module) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", filename, err)
	}

	c, err := ParseConfigData(data, modules)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var DecodeHookFuncs []mapstructure.DecodeHookFunc

func ParseConfigData(data []byte, modules []module.Module) (*Config, error) {
	var configMap map[string]interface{}
	err := yaml.Unmarshal(data, &configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	var c Config
	c.Modules = make(module.Modules, len(modules))

	var decoder *mapstructure.Decoder
	decoder, err = configDecoder(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to create config decoder: %w", err)
	}

	err = decoder.Decode(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config map: %w", err)
	}

	for _, m := range modules {
		c.Modules[m.Name()] = m
	}
	err = c.decodeModules(configMap)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) decodeModules(configMap map[string]interface{}) error {
	for k, v := range c.Modules {
		if m, ok := configMap[k]; ok {
			v := v
			decoder, err := configDecoder(&v)

			if err != nil {
				return fmt.Errorf("failed to create decoder for module %q: %w", k, err)
			}
			err = decoder.Decode(m)

			if err != nil {
				return fmt.Errorf("failed to decode module %q: %w", k, err)
			}
			c.Modules[k] = v
		}
	}
	return nil
}

func configDecoder(results interface{}) (*mapstructure.Decoder, error) {
	DecodeHookFuncs = append(DecodeHookFuncs, stringToTimeHookFunc())

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
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
			DecodeHookFuncs...,
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create mapstructure decoder: %w", err)
	}
	return decoder, nil
}

// stringToTimeHookFunc decodes either a simple time as am/pm or a RFC3339 formated time.
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

		dataStr, ok := data.(string)
		if !ok {
			return data, nil
		}

		if strings.HasSuffix(dataStr, "pm") || strings.HasSuffix(dataStr, "am") {
			parsedTime, err := time.Parse("15:04pm", dataStr)
			if err != nil {
				return data, fmt.Errorf("failed to parse time with format 15:04pm: %w", err)
			}
			return parsedTime, nil
		}

		if strings.HasSuffix(dataStr, "PM") || strings.HasSuffix(dataStr, "AM") {
			parsedTime, err := time.Parse("15:04PM", dataStr)
			if err != nil {
				return data, fmt.Errorf("failed to parse time with format 15:04PM: %w", err)
			}
			return parsedTime, nil
		}

		parsedTime, err := time.Parse(time.RFC3339, dataStr)
		if err != nil {
			return data, fmt.Errorf("failed to parse time with RFC3339 format: %w", err)
		}
		return parsedTime, nil
	}
}
