package main

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/hass-ws/services"
	"go.uber.org/zap"
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

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("zap logging could not be initialized: %v", err))
	}

	gs, err := goscript.New(config, zapr.NewLogger(zapLog))
	if err != nil {
		panic(err)
	}

	gs.AddTrigger(&goscript.Trigger{
		Unique: &goscript.Unique{KillMe: true},
		//Triggers:      []string{"input_button.test_button"},
		//DomainTrigger: []string{"input_button"},
		Periodic: goscript.Periodics(""),
		//Periodic: goscript.Periodics("*/1 * * * *"),
		States: goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
		Eval:   nil,
		Func: func(t *goscript.Task) {
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
			t.Sleep(10 * time.Second)
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
		},
	})
	gs.AddTrigger(&goscript.Trigger{
		Unique: &goscript.Unique{KillMe: true},
		//Triggers:      []string{"input_button.test_button"},
		//DomainTrigger: []string{"input_button"},
		//Periodic: goscript.Periodics(""),
		Periodic: goscript.Periodics("*/10 * * * * *"),
		States:   goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
		Eval:     nil,
		Func: func(t *goscript.Task) {
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
			t.Sleep(10 * time.Second)
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
