package control

import (
	"github.com/kjbreil/goscript/pkg/history"
	"github.com/kjbreil/goscript/pkg/state"
	"github.com/kjbreil/hass-ws/model"
)

type Responses chan Response

type Response struct {
	To   string
	From string

	Module     any
	States     *state.States
	ServiceRsp *model.Message
	Histories  *history.Histories
}
