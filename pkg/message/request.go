package message

import "github.com/kjbreil/goscript/pkg/trigger"

type RequestChan chan Request

type Request struct {
	To   string
	From string

	Trigger *trigger.Trigger

	// GobType []byte
}
