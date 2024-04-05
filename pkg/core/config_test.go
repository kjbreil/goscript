package core

import (
	"github.com/kjbreil/goscript/pkg/module"
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	type temperature struct {
		Fahrenheit float64
	}
	type highLow struct {
		High *temperature
		Loq  *temperature
	}
	type Lights struct {
		Name      string
		SnakeCase string
		Entities  []string
		OnTime    time.Time
		TestArray []highLow
	}

	var m module.Modules = make(module.Modules)

	m["lights"] = &Lights{}

	c, err := ParseConfig("config.yml", modules)
	if err != nil {
		t.Fatal(err)
	}
	if c.Websocket == nil {
		t.Fatal("Websocket is nil")
	}
	if c.MQTT == nil {
		t.Fatal("MQTT is nil")
	}
	l, ok := c.Modules["lights"]
	if !ok {
		t.Fatal("Lights is nil")
	}
	if l == nil {
		t.Fatal("Lights is nil")
	}
	if l.(*Lights).Name != "test" {
		t.Fatal("Lights.Name is not test")
	}
	if l.(*Lights).SnakeCase != "yes" {
		t.Fatal("Lights.SnakeCase is not yes")
	}
}
