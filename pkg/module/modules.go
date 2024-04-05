package module

type Modules map[string]Module

func (m Modules) Get(name string) (Module, bool) {
	if mo, ok := m[name]; ok {
		return mo, true
	}
	return nil, false
}
