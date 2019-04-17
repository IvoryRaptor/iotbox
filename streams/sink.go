package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Sink struct {
	Flow
	work func(msg interface{})
}

func (s *Sink) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		s.start(context)
	case interface{}:
		s.work(msg)
	}
}
