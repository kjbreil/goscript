package main

import (
	"fmt"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/hass-ws/model"
	"log"
	"os"
	"os/signal"
	"syscall"
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
		EntityTriggers: []string{"input_button.test_button"},
		EntityStates:   []string{"input_boolean.test_toggle"},
		Func: func(message model.Message, states []goscript.State) {
			fmt.Println(message.EntityID())
		},
	})

	gs.Connect()

	gs.GetStates([]string{"input_button.test_button", "input_boolean.test_toggle"})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Everything is set up")
	<-done

	gs.Close()
}
