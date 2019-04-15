package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/akka"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type TestTask struct {
	index int
}

func (t *TestTask) Init() {
	t.index = 0
	return
}

func (t *TestTask) GetNext(response *common.Response) *common.Request {
	var result *common.Request
	switch response.State {
	case common.Timeout:

	default:
		switch t.index {
		case 0:
			result = &common.Request{
				Wait: 1 * time.Second,
				Body: common.Message{"name": "a"},
			}
		case 1:
			result = &common.Request{
				Wait: 1 * time.Second,
				Body: common.Message{"name": "b"},
			}

		case 2:
			result = &common.Request{
				Wait: 1 * time.Second,
				Body: common.Message{"name": "c"},
			}
		}
		t.index++
	}
	return result
}

type Protocol interface {
}

type Port struct {
	protocol Protocol
}

func (port *Port) Receive(context akka.Context) {
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
