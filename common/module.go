package common

import (
	"github.com/IvoryRaptor/iotbox/akka"
	"time"
)

type Module struct {
	Port *akka.ActorRef
}

func (module *Module) Receive(context akka.Context) {
	switch task := context.Message().(type) {
	case ITask:
		var request Message = nil
		var response Message
		var future *akka.Future
		request = task.Init()
		for request != nil {
			future = context.Ask(module.Port, request, 1*time.Second)
			if result, err := future.Result(); err != nil {
				println(err.Error())
			} else {
				response = result.(Message)
				request = task.GetNext(response)
			}
		}
	}
}
