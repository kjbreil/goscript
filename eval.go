package goscript

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/kjbreil/hass-ws/model"
	"strconv"
)

func Eval(exp ...string) []string {
	return exp
}

func Evaluates(states States, eval []string) bool {
	passed := false
	for _, e := range eval {
		if Evaluate(states, e) {
			passed = true
		}
	}
	return passed
}

func Evaluate(states States, eval string) bool {
	var passed bool

	atoi := expr.Function(
		"float",
		func(params ...any) (any, error) {
			return strconv.ParseFloat(params[0].(string), 64)
		},
	)

	program, err := expr.Compile(eval, expr.Env(map[string]interface{}{}),
		expr.AllowUndefinedVariables(),
		expr.AsBool(),
		atoi)
	if err != nil {
		return false
	}

	env := make(map[string]interface{})

	if len(states) == 1 {
		for _, state := range states {
			env["state"] = state.State
			// add attributes to env
			if attr := state.Attributes; attr != nil {
				for k, v := range attr {
					for _, c := range program.Constants {
						if _, ok := c.(string); ok {
							if c.(string) == k {
								env[c.(string)] = v
							}
						}
					}
				}
			}
		}
	}

	for _, state := range states {
		env[state.DomainEntity] = state.State
		if attr := state.Attributes; attr != nil {
			for k, v := range attr {
				for _, c := range program.Constants {
					switch c.(type) {
					case string:
						if k == c.(string) {
							env[fmt.Sprintf("%s.%s", state.DomainEntity, c.(string))] = v
						}
					}
				}
			}
		}
	}

	evald, err := expr.Run(program, env)

	if err != nil {
		// TODO: Add error to some display
		return false
	}
	if evald.(bool) && !passed {
		passed = true
	}

	return passed
}

func (t *Trigger) eval(message *model.Message) bool {
	passed := !(len(t.Eval) > 0)

	states := map[string]*State{
		message.DomainEntity(): MessageState(message),
	}
	for _, e := range t.Eval {
		if Evaluate(states, e) {
			passed = true
		}
	}
	return passed
}
