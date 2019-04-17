package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Source struct {
	self *akka.ActorRef
	Flow
}

func (s *Source) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		s.start(context)
	case interface{}:
		s.tellNext(msg)
	}
}

func (s *Source) Start() {
	s.self = akka.EmptyRootContext.ActorOf(akka.PropsFromProducer(func() akka.Actor {
		return s
	}))
}

func (s *Source) Write(message interface{}) {
	akka.EmptyRootContext.Tell(s.self, message)
}
