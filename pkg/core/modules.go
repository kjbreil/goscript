package core

import "github.com/kjbreil/goscript/pkg/module"

func (gs *Core) ModuleMap() module.Modules {
	return gs.config.Modules
}

func (gs *Core) UpdateModule(key string, m module.Module) {
	gs.config.Modules[key] = m
}
