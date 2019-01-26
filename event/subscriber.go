package event

import (
	"github.com/ugo/podgo/handler"
)

type Subscriber struct {
	Handler handler.Handler
}

func (s *Subscriber) On(eventName string) error {
	return s.Handler.Handle(eventName)
}
