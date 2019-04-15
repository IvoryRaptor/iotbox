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
		var request *Request = nil
		var response Message
		var future *akka.Future
		request = task.GetNext(&Response{State: Initialize})
		for request != nil {
			future = context.Ask(module.portRef, request.Body, request.Wait)
			if result, err := future.Result(); err != nil {
				request = task.GetNext(&Response{
					State: Timeout,
					Body:  nil,
				})
			} else {
				response = result.(Message)
				request = task.GetNext(&Response{
					State: Result,
					Body:  response,
				})
			}
		}
	}
}
