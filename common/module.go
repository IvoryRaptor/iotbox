package common

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Module struct {
	Port *akka.ActorRef
}

func (module *Module) Receive(context akka.Context) {
	switch task := context.Message().(type) {
	case ITask:
		var request *Request = nil
		var response Message
		var future *akka.Future
		request = task.GetNext(&Response{State: Initialize})
		for request != nil {
			future = context.Ask(module.Port, request.Body, request.Wait)
			if result, err := future.Result(); err != nil {
				println(err.Error())
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
