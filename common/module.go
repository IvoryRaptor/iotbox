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
	case *TaskRef:
		request := task.GetRequest()
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
			request = task.GetRequest()
		}
	}
}

func CreatePort(port *Port, name string) (*akka.ActorRef, error) {
	return akka.EmptyRootContext.ActorOfNamed(akka.PropsFromProducer(func() akka.Actor {
		return &Module{Port: port}
	}), name)
}

func CreateTask(task ITask) *TaskRef {
	result := &TaskRef{
		task:            task,
		data:            map[string]interface{}{},
		func_receive:    task.Receive,
		func_getrequest: task.GetRequest,
	}
	task.Init(result)
	return result
}
