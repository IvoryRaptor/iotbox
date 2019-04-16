package common

import (
	"github.com/IvoryRaptor/iotbox/akka"
	"time"
)

func CreateActivePort(name string, port Port, wait time.Duration, work func(message Message)) (*akka.ActorRef, error) {
	result, err := akka.EmptyRootContext.ActorOfNamed(akka.PropsFromProducer(func() akka.Actor {
		return &ActiveRef{
			Port: port,
			Wait: wait,
			Work: work,
		}
	}), name)
	if err == nil {
		akka.EmptyRootContext.Tell(result, &Idle{})
	}
	return result, err
}

func CreatePassivePort(name string, port Port) (*akka.ActorRef, error) {
	return akka.EmptyRootContext.ActorOfNamed(akka.PropsFromProducer(func() akka.Actor {
		return &PassiveRef{Port: port}
	}), name)
}

func AddTask(module string, task ITask) *TaskRef {
	result := &TaskRef{
		task:            task,
		data:            map[string]interface{}{},
		func_receive:    task.Receive,
		func_getrequest: task.GetRequest,
	}
	task.Init(result)
	result.JoinProcessor(module)
	return result
}
