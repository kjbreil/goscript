// Package homekit provides HomeKit accessory protocol (HAP) integration for GoScript.
package homekit

import (
	"context"
	"fmt"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

// HomeKitConfig contains configuration for the HomeKit bridge.
type HomeKitConfig struct {
	BridgeName string
	Pin        string
}

// HomeKit manages a HomeKit accessory server with lifecycle control.
type HomeKit struct {
	ctx     context.Context
	cancel  context.CancelFunc
	mainCtx context.Context
}

// New creates a new HomeKit instance with the given main context.
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
			return fmt.Errorf("main context error: %w", h.mainCtx.Err())
		}
		h.cancel()
	}
	h.ctx, h.cancel = context.WithCancel(h.mainCtx)
	fs := hap.NewFsStore("./db")

	//nolint:exhaustruct // Only Name required; other Info fields optional
	a := accessory.NewBridge(accessory.Info{
		Name: "GoScript",
	})

	server, err := hap.NewServer(fs, a.A, accs...)
	if err != nil {
		return fmt.Errorf("failed to create HomeKit server: %w", err)
	}
	server.Pin = "00102003"

	err = server.ListenAndServe(h.ctx)
	if err != nil {
		return fmt.Errorf("failed to start HomeKit server: %w", err)
	}
	return nil
}
