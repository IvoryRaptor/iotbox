package common

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Module struct {
	Port    *Port
	portRef *akka.ActorRef
}

func (module *Module) Receive(context akka.Context) {
	switch task := context.Message().(type) {
	case *akka.Started:
		module.portRef = context.ActorOf(akka.PropsFromProducer(func() akka.Actor {
			return module.Port
		}))
	case ITask:
		request := task.GetNext()
		for request != nil {
			future := context.Ask(module.portRef, request.Body, request.Wait)
			if result, err := future.Result(); err != nil {
				task.Receive(&Response{
					State: Timeout,
					Body:  nil,
				})
			} else {
				task.Receive(&Response{
					State: Result,
					Body:  result.(Message),
				})
			}
			request = task.GetNext()
		}
	}
}

func CreatePort(port *Port, name string) (*akka.ActorRef, error) {
	return akka.EmptyRootContext.ActorOfNamed(akka.PropsFromProducer(func() akka.Actor {
		return &Module{Port: port}
	}), name)
}
