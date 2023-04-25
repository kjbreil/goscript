package main

import (
	"github.com/kjbreil/goscript"
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

	//gs.AddTrigger(&goscript.Trigger{
	//	Unique: &goscript.Unique{KillMe: true},
	//	//Triggers:      []string{"input_button.test_button"},
	//	//DomainTrigger: []string{"input_button"},
	//	Periodic: goscript.Periodics(""),
	//	//Periodic: goscript.Periodics("*/1 * * * *"),
	//	States: goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
	//	Eval:   nil,
	//	Func: func(t *goscript.Task) {
	//		gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
	//		t.Sleep(1 * time.Second)
	//		gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
	//	},
	//})
	gs.AddTrigger(&goscript.Trigger{
		//Unique: &goscript.Unique{KillMe: true},
		//Triggers:      []string{"input_button.test_button"},
		//DomainTrigger: []string{"input_button"},
		//Periodic: goscript.Periodics(""),
		Periodic: goscript.Periodics("*/3 * * * * *"),
		States:   goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
		Eval:     nil,
		Func: func(t *goscript.Task) {
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
			//t.Sleep(10 * time.Second)
			time.Sleep(10 * time.Second)
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
		},
	})

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
