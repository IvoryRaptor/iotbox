package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/akka"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type TestTask struct {
	index   int
	context *akka.RootContext
}

func (t *TestTask) Init(context *akka.RootContext) {
	t.index = 0
	t.context = context
	m := akka.NewLocalActorOf("com1")
	context.Tell(m, t)
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

func main() {
	rootContext := akka.EmptyRootContext

	rootContext.ActorOfNamed(akka.PropsFromProducer(func() akka.Actor {
		return &common.Module{Port: &common.Port{}}
	}), "com1")

	var task = &TestTask{}
	task.Init(rootContext)
	console.ReadLine()
}
