package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Source struct {
	self *akka.ActorRef
	BaseFlow
}

func (s *Source) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		s.Start(context)
	case interface{}:
		s.TellNext(msg)
	}
}

func (s *Source) Run() {
	s.self = akka.EmptyRootContext.ActorOf(akka.PropsFromProducer(func() akka.Actor {
		return s
	}))
}

func (s *Source) Write(message interface{}) {
	akka.EmptyRootContext.Tell(s.self, message)
}
