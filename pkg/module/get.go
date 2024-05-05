package module

import (
	"github.com/kjbreil/goscript/pkg/control"
	"github.com/kjbreil/goscript/pkg/history"
)

func GetHistories(m Module, histories history.GetHistories) {
	m.Requests().Chan() <- control.Request{
		To:      "goscript",
		From:    m.Name(),
		History: &histories,
	}
}
