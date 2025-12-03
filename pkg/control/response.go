package control

import (
	"github.com/kjbreil/goscript/pkg/history"
	"github.com/kjbreil/goscript/pkg/state"
	"github.com/kjbreil/hass-ws/model"
)

// Responses is a channel for sending Response messages.
type Responses chan Response

// Response represents a response message in the control system.
type Response struct {
	To   string
	From string

	Module     any
	States     *state.States
	ServiceRsp *model.Message
	Histories  *history.Histories
}
