package core

import (
	"github.com/google/uuid"
	"github.com/kjbreil/goscript/pkg/state"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/model"
)

func (gs *GoScript) handleGetStates(states []model.Result) {
	statesFuncToRun := make(map[uuid.UUID]*trigger.Task)

	for _, sr := range states {
		s := &state.State{
			DomainEntity: sr.DomainEntity(),
			Domain:       sr.Domain(),
			Entity:       sr.EntityID(),
			State:        state.StateText(sr.State()),
			Attributes:   sr.Attributes,
		}

		gs.states.Upsert(s)
	}

	for _, sr := range states {
		domainEntity := sr.DomainEntity()
		entityState := sr.State()
		message := &model.Message{
			Type: model.MessageTypeEvent,
			Event: &model.Event{
				Data: &model.Data{
					EntityId: &domainEntity,
					NewState: &model.State{
						EntityId:    &domainEntity,
						LastChanged: sr.LastChanged,
						State:       &entityState,
						Attributes:  sr.Attributes,
						LastUpdated: sr.LastUpdated,
						Context:     sr.Context,
					},
					OldState: nil,
				},
				EventType: model.EventTypeStateChanged,
				Context:   sr.Context,
			},
		}

		gs.Runner.RunTriggers(*message)
	}
	for _, t := range statesFuncToRun {
		gs.Runner.AddTask(t)
	}
}

func (gs *GoScript) handleHassMessage(message model.Message) {
	if message.Type == model.MessageTypeEvent {
		switch message.Event.EventType {
		case model.EventTypeStateChanged:

			s := &state.State{
				DomainEntity: message.DomainEntity(),
				Domain:       message.Domain(),
				Entity:       message.EntityID(),
				State:        state.StateText(message.State()),
				Attributes:   message.Attributes(),
			}

			gs.states.Upsert(s)

			gs.Runner.RunTriggers(message)
		case model.EventTypeCallService:

			gs.Runner.RunServiceTriggers(message)
		}
	}
}
