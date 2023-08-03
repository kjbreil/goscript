package main

import (
	"fmt"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
	"github.com/kjbreil/hass-ws/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := goscript.ParseConfig("config.yml", nil)
	if err != nil {
		panic(err)
	}

	gs, err := goscript.New(config, goscript.DefaultLogger())
	if err != nil {
		panic(err)
	}

	gs.AddTrigger(&goscript.Trigger{
		Unique: &goscript.Unique{KillMe: false},
		//Triggers:      []string{"input_button.test_button"},
		//DomainTrigger: []string{"input_button"},
		//Periodic: goscript.Periodics(""),
		Periodic: goscript.Periodics("*/3 * * * * *"),
		States:   goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
		Eval:     nil,
		Func: func(t *goscript.Task) {
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
			time.Sleep(10 * time.Second)
			//gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
		},
	})

	mainDevice := device.New("Task Switch", "task_switch", "Switch 1000", "Kaygel", "0.0.1")

	d, err := gs.AddDevice(mainDevice)
	if err != nil {
		panic(err)
	}

	switchOptions := entities.NewSwitchOptions()
	switchOptions.Name("Task Switch").
		CommandFunc(gs.TaskMQTT(&goscript.Trigger{
			States: []string{"input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"},
			Unique: &goscript.Unique{},
			Func: func(t *goscript.Task) {
				now := time.Now()
				gs.Logger().Info(fmt.Sprintf("Task Toggled - Bad Sleep - %s", now.Format(time.RFC3339)))
				t.Sleep(10 * time.Second)
				gs.Logger().Info(fmt.Sprintf("After Sleep - %s", now.Format(time.RFC3339)))
			},
		}))

	switchDevice, err := entities.NewSwitch(switchOptions)

	if err != nil {
		panic(err)
	}

	err = d.AddEntities([]entities.Entity{switchDevice})
	if err != nil {
		panic(err)
	}

	err = gs.Connect()
	if err != nil {
		panic(err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	gs.Logger().Info("Everything is set up")
	<-done

	gs.Close()
}
