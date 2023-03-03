package goscript

import (
	"context"
	"github.com/antonmedv/expr"
	"github.com/google/uuid"
	"github.com/kjbreil/hass-ws/model"
	"strconv"
)

type Trigger struct {
	uuid uuid.UUID

	Unique   *Unique
	Triggers []string
	States   []string
	Eval     []string
	Func     TriggerFunc
}

type Unique struct {
	KillMe bool // even when false is true
	ctx    context.Context
	cancel context.CancelFunc
}

type TriggerFunc func(t *Task, message model.Message, states States)

func (gs *GoScript) AddTrigger(t *Trigger) {
	t.uuid = uuid.New()
	if t.Unique != nil {
		t.Unique.ctx, t.Unique.cancel = context.WithCancel(context.Background())
	}
	entityTriggers := make(map[string]struct{})
	entityStates := make(map[string]struct{})

	for _, et := range t.Triggers {
		entityTriggers[et] = struct{}{}
		entityStates[et] = struct{}{}
	}
	for _, es := range t.States {
		entityStates[es] = struct{}{}
	}
	t.Triggers = nil
	t.States = nil
	t.Triggers = make([]string, len(entityTriggers))
	i := 0
	for k := range entityTriggers {
		t.Triggers[i] = k
		i++
	}

	i = 0
	t.States = make([]string, len(entityStates))
	for k := range entityStates {
		t.States[i] = k
		i++
	}

	for _, et := range t.Triggers {
		gs.triggers[et] = append(gs.triggers[et], t)
	}
}

func (gs *GoScript) AddTriggers(triggers ...*Trigger) {
	for _, t := range triggers {
		gs.AddTrigger(t)
	}
}

func Eval(exp ...string) []string {
	return exp
}

func (gs *GoScript) runTriggers(message model.Message) {
	funcToRun := make(map[uuid.UUID]Task)

	if tr, ok := gs.triggers[message.DomainEntity()]; ok {
		for _, t := range tr {
			passed := !(len(t.Eval) > 0)
			for _, e := range t.Eval {
				atoi := expr.Function(
					"float",
					func(params ...any) (any, error) {
						return strconv.ParseFloat(params[0].(string), 64)
					},
				)

				program, err := expr.Compile(e, expr.Env(map[string]interface{}{}),
					expr.AllowUndefinedVariables(),
					expr.AsBool(),
					atoi)
				if err != nil {
					continue
				}

				env := make(map[string]interface{})
				env["state"] = message.State()

				// add attributes to env
				if attr := message.Attributes(); attr != nil {
					for k, v := range attr {
						for _, c := range program.Constants {
							switch c.(type) {
							case string:
								if k == c.(string) {
									env[c.(string)] = v
								}
							}

						}
					}
				}

				evald, err := expr.Run(program, env)

				if err != nil {
					// TODO: Add error to some display
					continue
				}
				if evald.(bool) && !passed {
					passed = true
				}
			}
			if passed {
				if t.Unique != nil {
					t.Unique.cancel()
					t.Unique.ctx, t.Unique.cancel = context.WithCancel(context.Background())
					funcToRun[t.uuid] = Task{
						f:      t.Func,
						states: t.States,
						ctx:    t.Unique.ctx,
						cancel: t.Unique.cancel,
					}
				} else {
					ctx, cancel := context.WithCancel(context.Background())
					funcToRun[uuid.New()] = Task{
						f:      t.Func,
						states: t.States,
						ctx:    ctx,
						cancel: cancel,
					}
				}
			}
		}
	}

	for _, t := range funcToRun {
		ss := gs.GetStates(t.states)
		go t.f(&t, message, ss)
	}
}
