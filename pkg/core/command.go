package core

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
	//nolint:exhaustruct // Only service type fields are set; entities unused
	return &Command{
		t:       CommandTypeService,
		service: service,
	}
}

func GetStatesCommand(entities ...string) *Command {
	//nolint:exhaustruct // Only getStates type fields are set; service unused
	return &Command{
		t:        CommandTypeGetStates,
		entities: entities,
	}
}
