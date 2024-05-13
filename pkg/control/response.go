package control

import "github.com/kjbreil/goscript/pkg/state"

type Responses chan Response

type Response struct {
	To   string
	From string

	Module any
	States *state.States
}
