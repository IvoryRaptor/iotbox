package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/akka"
	"github.com/IvoryRaptor/iotbox/common"
)

type TestTask struct {
	index int
}

func (t *TestTask) Init() common.Message {
	t.index = 1
	return common.Message{"name": "a"}
}

func (t *TestTask) GetNext(msg common.Message) common.Message {
	var result common.Message
	switch t.index {
	case 1:
		result = common.Message{"name": "b"}
	case 2:
		result = common.Message{"name": "c"}
	}
	t.index++
	return result
}

type Protocol interface {
}

type Port struct {
	protocol Protocol
}

func (module *Port) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case common.Message:
		println(msg["name"].(string))
		response := common.Message{
			"name":  msg["name"].(string),
			"value": "1",
		}
		context.Tell(context.Sender(), response)
	}
}

func main() {
	rootContext := akka.EmptyRootContext
	port := rootContext.ActorOf(akka.PropsFromProducer(func() akka.Actor {
		return &Port{}
	}))

	tmp := rootContext.ActorOf(akka.PropsFromProducer(func() akka.Actor {
		return &common.Module{Port: port}
	}))

	rootContext.Tell(tmp, &TestTask{})

	console.ReadLine()
}
