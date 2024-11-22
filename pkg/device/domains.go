//go:generate stringer -type=DomainType -trimprefix=DomainType
//go:generate jsonenums -type=DomainType

package device

import "sync"

type Domains struct {
	d  map[DomainType]*Domain
	mu *sync.Mutex
}

func NewDomains() *Domains {
	return &Domains{
		d:  make(map[DomainType]*Domain),
		mu: &sync.Mutex{},
	}
}

type DomainType int

const (
	DomainTypeUnknown DomainType = iota
	DomainTypeLight
	DomainTypeSwitch
	DomainTypeSensor
)

func (d *Domains) AddEntity(dt DomainType, e Entity) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.d[dt]; !ok {
		d.d[dt] = NewDomain(dt)
	}
	d.d[dt].AddEntity(e)
}

func (d *Domains) GetEntity(dt DomainType, entityName string) Entity {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.d[dt]; !ok {
		return nil
	}
	return d.d[dt].GetEntity(entityName)
}
