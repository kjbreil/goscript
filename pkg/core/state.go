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
		//nolint:exhaustruct // LastChanged and LastUpdated not needed for initial state setup
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
		//nolint:exhaustruct // Message fields set as needed for state change event
		message := &model.Message{
			Type: model.MessageTypeEvent,
			Event: &model.Event{ //nolint:exhaustruct // Event fields set for state change
				Data: &model.Data{ //nolint:exhaustruct // Data fields set for state change
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

			//nolint:exhaustruct // LastChanged and LastUpdated not needed for state updates
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
		case model.EventTypeAll:
			// EventTypeAll is a subscription type, not an actual event
			// No action needed for this case
		case model.EventTypeState:
			// EventTypeState is used for state subscriptions
			// Actual state events use EventTypeStateChanged
		case model.EventTypePysScriptRunning:
			// PysScript running events - not currently handled
			// Could be logged or monitored if needed in the future
		}
	}
}
