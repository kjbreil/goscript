package core

import "github.com/kjbreil/goscript/pkg/module"

func (gs *GoScript) ModuleMap() module.Modules {
	return gs.config.Modules
}

func (gs *GoScript) UpdateModule(key string, m module.Module) {
	gs.config.Modules[key] = m
}
