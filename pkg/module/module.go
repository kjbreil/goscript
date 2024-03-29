package module

type Module interface {
	Update() error
	Name() string
	Close() error

	mustImplementBase()
}

type Base struct {
}

func (m *Base) mustImplementBase() {}
