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

	Unique        *Unique
	Triggers      []string
	DomainTrigger []string // DomainTrigger, triggers of everything in the domain, also attaches all states for the domain
	Periodic
	DomainStates []string
	States       []string
	Eval         []string
	Func         TriggerFunc
}

type Unique struct {
	KillMe bool // even when false is true
	ctx    context.Context
	cancel context.CancelFunc
}

type TriggerFunc func(t *Task)

func Entities(entities ...string) []string {
	var rtn []string
	for _, e := range entities {
		rtn = append(rtn, e)
	}
	return rtn
}
func (gs *GoScript) AddTrigger(t *Trigger) {
	// setup the trigger object
	t.uuid = uuid.New()
	if t.Unique != nil {
		t.Unique.ctx, t.Unique.cancel = context.WithCancel(context.Background())
	}
	entityTriggers := make(map[string]struct{})
	domainTriggers := make(map[string]struct{})
	entityStates := make(map[string]struct{})
	domainStates := make(map[string]struct{})
	for _, et := range t.Triggers {
		entityTriggers[et] = struct{}{}
		entityStates[et] = struct{}{}
	}
	for _, es := range t.States {
		entityStates[es] = struct{}{}
	}
	for _, ed := range t.DomainTrigger {
		domainTriggers[ed] = struct{}{}
	}
	for _, eds := range t.DomainStates {
		domainStates[eds] = struct{}{}
	}

	t.Triggers = nil
	t.States = nil
	t.DomainTrigger = nil
	t.DomainStates = nil

	t.Triggers = make([]string, len(entityTriggers))
	i := 0
	for k := range entityTriggers {
		t.Triggers[i] = k
		i++
	}

	t.DomainTrigger = make([]string, len(domainTriggers))
	i = 0
	for k := range domainTriggers {
		t.DomainTrigger[i] = k
		i++
	}

	t.States = make([]string, len(entityStates))
	i = 0
	for k := range entityStates {
		t.States[i] = k
		i++
	}

	t.DomainStates = make([]string, len(domainStates))
	i = 0
	for k := range domainStates {
		t.DomainStates[i] = k
		i++
	}

	// for each entity add to the triggers map
	for _, et := range t.Triggers {
		gs.triggers[et] = append(gs.triggers[et], t)
	}

	// for each domain add to the domain trigger map
	for _, edt := range t.DomainTrigger {
		gs.domainTrigger[edt] = append(gs.domainTrigger[edt], t)
	}

	// for each periodic add to the periodic map
	// cron time is an array of triggers so multiple triggers can have same cron schedule
	for _, ep := range t.Periodic {
		gs.periodic[ep] = append(gs.periodic[ep], t)
	}
}

func (gs *GoScript) RemoveTrigger(t *Trigger) {
	for _, et := range t.Triggers {
		for i, te := range gs.triggers[et] {
			if te.uuid == t.uuid {
				gs.triggers[et] = append(gs.triggers[et][:i], gs.triggers[et][i+1:]...)
				break
			}
		}
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

	if tr, ok := gs.triggers[message.DomainEntity()]; ok {
		for _, trigger := range tr {
			gs.triggerDomainEntity(&message, trigger)
		}
	}

	if tr, ok := gs.domainTrigger[message.Domain()]; ok {
		for _, trigger := range tr {
			gs.triggerDomain(&message, trigger)
		}
	}

}

func (gs *GoScript) triggerDomainEntity(message *model.Message, trigger *Trigger) {
	passed := trigger.eval(message)
	if passed {
		task := gs.newTask(trigger, message)
		gs.funcToRun[task.uuid] = task
	}
}
func (gs *GoScript) triggerDomain(message *model.Message, trigger *Trigger) {
	passed := trigger.eval(message)
	if passed {
		task := gs.newTask(trigger, message)
		gs.funcToRun[task.uuid] = task
	}
}

func (tr *Trigger) eval(message *model.Message) bool {
	passed := !(len(tr.Eval) > 0)
	for _, e := range tr.Eval {
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

	return passed

}
