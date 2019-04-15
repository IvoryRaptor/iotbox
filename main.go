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
	m := akka.NewLocalActorOf("com1")
	akka.EmptyRootContext.Tell(m, t)
	return
}

func (t *TestTask) Receive(response *common.Response) {
	switch response.State {
	case common.Timeout:

	default:
	}
	t.index++
}

func (t *TestTask) GetNext() *common.Request {
	switch t.index {
	case 0:
		return &common.Request{
			Wait: 1 * time.Second,
			Body: common.Message{"name": "a"},
		}
	case 1:
		return &common.Request{
			Wait: 1 * time.Second,
			Body: common.Message{"name": "b"},
		}

	case 2:
		return &common.Request{
			Wait: 1 * time.Second,
			Body: common.Message{"name": "c"},
		}
	default:
		return nil
	}
}

func main() {
	common.CreatePort(&common.Port{}, "com1")
	test := TestTask{}
	test.Init()
	console.ReadLine()
}
