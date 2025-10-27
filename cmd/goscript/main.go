package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kjbreil/goscript/modules/circadian"
	"github.com/kjbreil/goscript/modules/lights"
	"github.com/kjbreil/goscript/modules/virtual"
	"github.com/kjbreil/goscript/pkg/core"
	"github.com/kjbreil/goscript/pkg/module"
)

func main() {
	ms := []module.Module{
		&lights.Lights{},
		&virtual.Virtual{},
		&circadian.Circadian{},
	}
	config, err := core.ParseConfig("config.yml", ms)
	if err != nil {
		panic(err)
	}

	// Create logger with configured log level
	var logger *slog.Logger
	if config.GoScript != nil {
		logger = core.DefaultLoggerWithLevel(config.GoScript.GetLogLevel())
	} else {
		logger = core.DefaultLogger()
	}

	gs, err := core.New(config, logger)
	if err != nil {
		panic(err)
	}

	gs.UpdateModule("lights", core.GetModule[*lights.Lights](gs, "lights"))
	gs.UpdateModule("virtual", core.GetModule[*virtual.Virtual](gs, "virtual"))
	gs.UpdateModule("circadian", core.GetModule[*circadian.Circadian](gs, "circadian"))

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
