package homekit

import (
	"context"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

type HomeKitConfig struct {
	BridgeName string
	Pin        string
}

type HomeKit struct {
	ctx     context.Context
	cancel  context.CancelFunc
	mainCtx context.Context
}

func New(mainCtx context.Context) *HomeKit {
	ctx, cancel := context.WithCancel(mainCtx)
	return &HomeKit{
		mainCtx: mainCtx,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Run starts the HomeKit server or restarts the homekit server with new accessories.
func (h *HomeKit) Run(accs []*accessory.A) error {
	if h.ctx.Err() != nil {
		if h.mainCtx.Err() != nil {
			return h.mainCtx.Err()
		}
		h.cancel()
	}
	h.ctx, h.cancel = context.WithCancel(h.mainCtx)
	fs := hap.NewFsStore("./db")

	a := accessory.NewBridge(accessory.Info{
		Name: "GoScript",
	})

	server, err := hap.NewServer(fs, a.A, accs...)
	if err != nil {
		return err
	}
	server.Pin = "00102003"

	return server.ListenAndServe(h.ctx)
}
