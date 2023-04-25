package goscript

import (
	"fmt"
	"github.com/kjbreil/hass-ws/services"
)

// ServiceChan is a channel to send services to be run to
type ServiceChan chan services.Service

func (gs *GoScript) CallService(service services.Service) {
	gs.ws.CallService(service)
}

func (gs *GoScript) runService() {
	for {
		select {
		case <-gs.ctx.Done():
			return
		case s := <-gs.ServiceChan:
			gs.logger.V(4).Info(fmt.Sprintf("%s", s.JSON()))
			gs.ws.CallService(s)
		}
	}
}
