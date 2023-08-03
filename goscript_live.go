package goscript

import (
	"github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
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
	err = generateTestDevices(gs)
	if err != nil {
		t.Fatal(err)
	}

	for _, fn := range preFns {
		err = fn(gs)
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
		err = fn(gs)
		if err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()

	//done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//
	//<-done
}

type TestFunc func(gs *GoScript) error

func generateTestDevices(gs *GoScript) error {
	mainDevice := device.New("Test Devices", "test_devices", "Tester 1000", "goscript", "0.0.1")

	d, err := gs.AddDevice(mainDevice)
	if err != nil {
		panic(err)
	}

	switchOptions := entities.NewSwitchOptions()
	switchOptions.Name("Test Switch")

	switchDevice, err := entities.NewSwitch(switchOptions)

	if err != nil {
		return err
	}
	err = d.AddEntities([]entities.Entity{switchDevice})
	if err != nil {
		return err
	}

	return nil
}
