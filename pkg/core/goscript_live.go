package core

import (
	"github.com/kjbreil/goscript/pkg/control"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/device/entities"
	hassentities "github.com/kjbreil/hass-mqtt/entities"
	hassdevice "github.com/kjbreil/hass-mqtt/pkg/device"

	"sync"
	"testing"
)

func GoScriptTestRun(preFns []TestFunc, postFns []TestFunc, wg *sync.WaitGroup, t *testing.T) {
	config, err := ParseConfig("config.yml", nil)
	if err != nil {
		t.Fatal(err)
	}

	gs, err := New(config, DefaultLogger())
	if err != nil {
		t.Fatal(err)
	}

	// Generate Devices for testing
	d, err := generateTestDevices()
	if err != nil {
		t.Fatal(err)
	}

	gs.requests.Chan() <- control.Request{
		To:      "goscript",
		From:    "test",
		Trigger: nil,
		Device:  d,
	}

	for _, fn := range preFns {
		err = fn(gs.requests)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = gs.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer gs.Close()

	for _, fn := range postFns {
		err = fn(gs.requests)
		if err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()

	// done := make(chan os.Signal, 1)
	// signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//
	// <-done
}

type TestFunc func(ctrl *control.Requests) error

func generateTestDevices() (*device.Device, error) {
	mainDevice := hassdevice.New("Test Devices", "test_devices", "Tester 1000", "goscript", "0.0.1")

	d := device.NewDevice(mainDevice)

	switchOptions := hassentities.NewSwitchOptions()
	switchOptions.Name("Test Switch")

	switchDevice, err := hassentities.NewSwitch(switchOptions)
	if err != nil {
		return nil, err
	}
	err = d.AddEntities([]device.Entity{entities.NewHassEntity(switchDevice)})
	if err != nil {
		return nil, err
	}

	return d, nil
}
