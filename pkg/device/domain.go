package device

import "sync"

type Domain struct {
	e Entities
}

func NewDomain(dt DomainType) *Domain {
	return &Domain{
		e: Entities{
			e:  make(map[string]Entity),
			mu: &sync.Mutex{},
		},
	}
}

func (d *Domain) AddEntity(e Entity) {
	d.e.mu.Lock()
	defer d.e.mu.Unlock()
	d.e.e[e.GetName()] = e
}

func (d *Domain) GetEntity(name string) Entity {
	d.e.mu.Lock()
	defer d.e.mu.Unlock()
	return d.e.e[name]
}
