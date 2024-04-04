package main

import (
	"github.com/kjbreil/goscript/modules/motion"
	"github.com/kjbreil/goscript/pkg/core"
	"github.com/kjbreil/goscript/pkg/module"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ms := []module.Module{
		&motion.Motion{},
	}
	config, err := core.ParseConfig("config.yml", ms)
	if err != nil {
		panic(err)
	}

	gs, err := core.New(config, core.DefaultLogger())
	if err != nil {
		panic(err)
	}

	gs.UpdateModule("motion", core.GetModule[*motion.Motion](gs, "motion"))

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
