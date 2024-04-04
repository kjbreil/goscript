package state

func (s *States) GetDomainStates(domainentity []string) *States {
	rtn := NewMultiStates(s.FindDomainMap(domainentity))

	return &rtn
}
