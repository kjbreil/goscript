package goscript

import (
	"fmt"
	"github.com/kjbreil/hass-ws/model"
)

func (gs *GoScript) GetState(domain, entityid string) *State {
	txn := gs.db.Txn(false)
	defer txn.Abort()
	raw, err := txn.First("state", "id", fmt.Sprintf("%s.%s", domain, entityid))
	if err != nil {
		panic(err)
	}
	return raw.(*State)
}

func (gs *GoScript) GetStates(domainentity []string) []*State {
	var rtn []*State

	txn := gs.db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("state", "id", domainentity)
	if err != nil {
		panic(err)
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		s := obj.(*State)
		rtn = append(rtn, s)
	}
	return rtn
}

func (gs *GoScript) handleMessage(message model.Message) {
	switch message.Type {
	case model.MessageTypeEvent:

		switch message.Event.EventType {
		case model.EventTypeStateChanged:
			txn := gs.db.Txn(true)
			defer txn.Abort()

			s := &State{
				DomainEntity: message.DomainEntity(),
				Domain:       message.Domain(),
				Entity:       message.EntityID(),
				State:        message.State(),
				Attributes:   message.Attributes(),
			}

			if err := txn.Insert("state", s); err != nil {
				panic(err)
			}
			txn.Commit()

			gs.runTriggers(message)
		}
	}
}

func (gs *GoScript) handleGetStates(states []model.Result) {
	txn := gs.db.Txn(true)
	defer txn.Abort()

	for _, sr := range states {
		s := &State{
			DomainEntity: sr.DomainEntity(),
			Domain:       sr.Domain(),
			Entity:       sr.EntityID(),
			State:        sr.State(),
			Attributes:   sr.Attributes,
		}

		if err := txn.Insert("state", s); err != nil {
			panic(err)
		}
	}

	txn.Commit()
}
