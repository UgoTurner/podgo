package event

import (
	"github.com/ugo/podcastor/handler"
)

type Subscriber struct {
	Handler handler.Handler
}

func (s *Subscriber) On(eventName string) error {
	return s.Handler.Handle(eventName)
}
