package goscript

import "github.com/kjbreil/hass-ws/services"

func (gs *GoScript) CallService(service services.Service) {
	gs.ws.CallService(service)

}
