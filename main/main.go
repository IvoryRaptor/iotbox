package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type TestTask struct {
	index int
}

func (t *TestTask) Init(task *common.TaskRef) {
	t.index = 0
	return
}

func (t *TestTask) Receive(task *common.TaskRef, response *common.Response) {
	switch response.State {
	case common.Timeout:

	default:
	}
	t.index++
}

func (t *TestTask) GetRequest(task *common.TaskRef) *common.Request {
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
	common.JoinTask("com1", &TestTask{})
	console.ReadLine()
}