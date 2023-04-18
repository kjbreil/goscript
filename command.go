package goscript

import "github.com/kjbreil/hass-ws/services"

type CommandType int

const (
	CommandTypeService   CommandType = iota
	CommandTypeGetStates CommandType = iota
)

type Command struct {
	t        CommandType
	service  services.Service
	entities []string
}

func ServiceCommand(service services.Service) *Command {
	return &Command{
		t:       CommandTypeService,
		service: service,
	}
}

func GetStatesCommand(entities ...string) *Command {
	return &Command{
		t:        CommandTypeGetStates,
		entities: entities,
	}
}
