// Package eval provides expression evaluation functionality using the expr library.
package eval

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/kjbreil/goscript/pkg/state"
)

// Eval returns the provided expressions as-is for evaluation.
func Eval(exp ...string) []string {
	return exp
}

// Evaluates checks if any of the provided evaluation expressions pass for the given states.
func Evaluates(states state.States, eval []string) bool {
	passed := false
	for _, e := range eval {
		if Evaluate(states, e) {
			passed = true
		}
	}
	return passed
}

// Evaluate evaluates a single expression against the given states.
func Evaluate(states state.States, eval string) bool {
	var passed bool

	program, err := expr.Compile(eval, expr.Env(map[string]interface{}{}),
		expr.AllowUndefinedVariables(),
		expr.AsBool(),
	)
	if err != nil {
		return false
	}

	env := make(map[string]interface{})

	if states.Len() == 1 {
		for _, state := range states.Slice() {
			env["state"] = string(state.State)
			// add attributes to env
			if attr := state.Attributes; attr != nil {
				for k, v := range attr {
					for _, c := range program.Constants {
						if cStr, ok := c.(string); ok {
							if cStr == k {
								env[cStr] = v
							}
						}
					}
				}
			}
		}
	}

	for _, state := range states.Slice() {
		env[state.DomainEntity] = string(state.State)
		if attr := state.Attributes; attr != nil {
			for k, v := range attr {
				for _, c := range program.Constants {
					if cStr, ok := c.(string); ok {
						if k == cStr {
							env[fmt.Sprintf("%s.%s", state.DomainEntity, cStr)] = v
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
	if evaldBool, ok := evald.(bool); ok && evaldBool && !passed {
		passed = true
	}

	return passed
}
