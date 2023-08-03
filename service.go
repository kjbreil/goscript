package goscript

import (
	"fmt"
	"github.com/kjbreil/hass-ws/services"
	"strings"
	"time"
)

// ServiceChan is a channel to send services to be run to
type ServiceChan chan services.Service

func (gs *GoScript) CallService(service services.Service) {
	gs.ws.CallService(service)
}

func (gs *GoScript) runService() {
	chanBuffer := make(map[string]services.Service)
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-gs.ctx.Done():
			return
		case s := <-gs.ServiceChan:
			gs.logger.V(4).Info(fmt.Sprintf("%s", s.JSON()))
			entitiesString := strings.Join(s.Targets(), ",")
			chanBuffer[entitiesString] = s
		case <-ticker.C:
			for k, v := range chanBuffer {
				gs.ws.CallService(v)
				delete(chanBuffer, k)
			}
		}
	}
}
