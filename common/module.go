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
			retry := 0
			future := context.Ask(module.portRef, request.Body, request.Wait)
			for ; retry < request.Retry; retry++ {
				if result, err := future.Result(); err == nil {
					task.Receive(&Response{
						State: ResponseTimeout,
						Body:  result.(Message),
					})
					break
				}
			}
			if retry >= request.Retry {
				task.Receive(&Response{
					State: ResponseResult,
					Body:  nil,
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

func JoinTask(module string, task ITask) *TaskRef {
	result := &TaskRef{
		task:            task,
		data:            map[string]interface{}{},
		func_receive:    task.Receive,
		func_getrequest: task.GetRequest,
	}
	task.Init(result)
	result.JoinModule(module)
	return result
}
