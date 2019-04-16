package main

import (
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
	case common.ResponseTimeout:

	default:
	}
	t.index++
}

func (t *TestTask) GetRequest(task *common.TaskRef) *common.Request {
	switch t.index {
	case 0:
		return &common.Request{
			Wait:  1 * time.Second,
			Body:  common.Message{"name": "a"},
			Retry: 3,
		}
	case 1:
		return &common.Request{
			Wait:  1 * time.Second,
			Body:  common.Message{"name": "b"},
			Retry: 3,
		}

	case 2:
		return &common.Request{
			Wait:  1 * time.Second,
			Body:  common.Message{"name": "c"},
			Retry: 3,
		}
	default:
		return nil
	}
}
