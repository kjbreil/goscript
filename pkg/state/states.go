package state

func (ss *States) GetDomainStates(domainentity []string) *States {
	rtn := NewMultiStates(ss.FindDomainMap(domainentity))

	return &rtn
}

func (ss *States) GetState(entityID string) *State {
	s, _ := ss.Get(entityID)
	return s
}

func (ss *States) GetStates(domainentity []string) *States {
	rtn := NewMultiStates(ss.Find(domainentity))

	return &rtn
}
