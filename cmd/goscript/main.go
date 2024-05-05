package main

import (
	"github.com/kjbreil/goscript/modules/circadian"
	"github.com/kjbreil/goscript/modules/lights"
	"github.com/kjbreil/goscript/modules/virtual"
	"github.com/kjbreil/goscript/pkg/core"
	"github.com/kjbreil/goscript/pkg/module"
	"os"
	"os/signal"
	"syscall"
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

	gs, err := core.New(config, core.DefaultLogger())
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
