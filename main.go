package main

import (
	"fmt"
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/akka"
	"time"
)

type Message map[string]interface{}

type TaskStart struct {
}

type Task struct {
	result map[string]interface{}
}

func (task *Task) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		task.result = map[string]interface{}{}
		child := akka.NewLocalActorOf("com1")
		context.Tell(child, context.Self())
		println("Task Started")
	case TaskStart:
		println("task a")
		context.Tell(context.Sender(), Message{"name": "a"})
	case Message:
		switch msg["name"] {
		case "a":
			task.result["a"] = msg["value"]
			context.Tell(context.Sender(), Message{"name": "b"})
			println("task b")
		case "b":
			task.result["b"] = msg["value"]
			context.Tell(context.Sender(), Message{"name": "c"})
			println("task c")
		case "c":
			task.result["c"] = msg["value"]
			fmt.Printf("a = %s\n b = %s\n c = %s\n",
				task.result["a"],
				task.result["b"],
				task.result["c"],
			)
			context.Tell(context.Sender(), nil)
		}
	}
}

type Protocol interface {
}

type Port struct {
	protocol Protocol
}

func (module *Port) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case Message:
		response := Message{
			"name":  msg["name"].(string),
			"value": "1",
		}
		context.Tell(context.Sender(), response)
	}
}

type Module struct {
	Port *akka.ActorRef
}

func (module *Module) Receive(context akka.Context) {
	switch task := context.Message().(type) {
	case *akka.ActorRef:
		var request Message = nil
		var response interface{} = TaskStart{}
		var future *akka.Future
		for ok := true; ok; ok = request != nil {
			future = context.Ask(task, response, 1*time.Second)
			if result, err := future.Result(); err != nil {
				println(err.Error())
				request = nil
			} else if result != nil {
				request = result.(Message)
				future = context.Ask(module.Port, request, 1*time.Second)
				if result, err := future.Result(); err != nil {
					println(err.Error())
				} else {
					response = result.(Message)
				}
			} else {
				request = nil
			}
		}
	}
}

func main() {
	rootContext := akka.EmptyRootContext
	port := rootContext.ActorOf(akka.PropsFromProducer(func() akka.Actor {
		return &Port{}
	}))
	rootContext.ActorOfNamed(akka.PropsFromProducer(func() akka.Actor {
		return &Module{Port: port}
	}), "com1")
	rootContext.ActorOf(akka.PropsFromProducer(func() akka.Actor {
		return &Task{}
	}))
	console.ReadLine()
}
