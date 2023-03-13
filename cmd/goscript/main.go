package main

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/hass-ws/services"
	"go.uber.org/zap"
	"log"
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

	gs, err := goscript.New(config)
	if err != nil {
		panic(err)
	}

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("zap logging could not be initialized", err))
	}
	gs.Logger(zapr.NewLogger(zapLog))

	gs.AddTrigger(&goscript.Trigger{
		Unique: &goscript.Unique{KillMe: true},
		//Triggers: []string{"input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"},
		Triggers: []string{"input_button.test_button"},
		Periodic: goscript.Periodics("* * * * *", ""),
		States:   goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
		Eval:     nil,
		//Eval:     goscript.Eval(`float(state) > 10`, `state == "on"`),
		Func: func(t *goscript.Task) {
			log.Println("flipping")
			log.Println(t.States["input_boolean.test_toggle"].State)
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
			log.Println("waiting 5 seconds")

			t.Sleep(10 * time.Second)

			//if ok := t.WaitUntil("input_number.test_number", nil, 5*time.Second); !ok {
			//	return
			//}
			log.Println("flipping at end")
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
			log.Println(t.States["input_boolean.test_toggle"].State)

		},
	})

	err = gs.Connect()
	if err != nil {
		panic(err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Everything is set up")
	<-done

	gs.Close()
}
