package main

import (
	"fmt"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/hass-ws/model"
	"github.com/kjbreil/hass-ws/services"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := goscript.ParseConfig("config.yml")
	if err != nil {
		panic(err)
	}

	gs, err := goscript.New(config)
	if err != nil {
		panic(err)
	}

	gs.AddTrigger(&goscript.Trigger{
		//Triggers: []string{"input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"},
		Triggers: []string{"input_button.test_button"},
		States:   []string{},
		Unique:   &goscript.Unique{KillMe: true},
		//Eval:     goscript.Eval(`float(state) > 10`, `state == "on"`),
		Func: func(t *goscript.Task, message model.Message, states goscript.States) {
			fmt.Println("flipping")
			gs.CallService(services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle")))
			fmt.Println("waiting 5 seconds")
			if ok := t.Sleep(5 * time.Second); !ok {
				return
			}
			fmt.Println("flipping at end")
			gs.CallService(services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle")))

		},
	})

	gs.Connect()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Everything is set up")
	<-done

	gs.Close()
}
