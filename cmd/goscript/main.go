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
		panic(fmt.Sprintf("zap logging could not be initialized", err))
	}

	gs, err := goscript.New(config, zapr.NewLogger(zapLog))
	if err != nil {
		panic(err)
	}

	gs.AddTrigger(&goscript.Trigger{
		Unique:        &goscript.Unique{KillMe: true},
		Triggers:      []string{"input_button.test_button"},
		DomainTrigger: []string{"input_button"},
		Periodic:      goscript.Periodics("* * * * *", ""),
		States:        goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
		Eval:          nil,
		Func: func(t *goscript.Task) {
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
			t.Sleep(10 * time.Second)
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
		},
	})

	err = gs.Connect()

	if err != nil {
		panic(err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	gs.GetLogger().Info("Everything is set up")
	<-done

	gs.Close()
}
