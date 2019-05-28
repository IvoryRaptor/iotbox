package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Sink struct {
	BaseFlow
	work func(msg interface{})
}

func (s *Sink) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		s.Start(context)
	case interface{}:
		s.work(msg)
	}
}
