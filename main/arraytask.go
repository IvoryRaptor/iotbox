package main

import (
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type ArrayTask struct {
	index      int
	Messages   []common.Message
	Wait       time.Duration
	retryCount int
	Retry      int
}

func (t *ArrayTask) Init(task *common.TaskRef) {
	t.index = 0
	t.retryCount = 0
	return
}

func (t *ArrayTask) Receive(task *common.TaskRef, response *common.Response) {
	switch response.State {
	case common.ResponseTimeout:
		t.retryCount++
		if t.retryCount >= t.Retry {
			t.index++
		}
	default:
		t.retryCount = 0
		t.index++
	}
}

func (t *ArrayTask) GetRequest(task *common.TaskRef) *common.Request {
	if t.index >= len(t.Messages) {
		return nil
	}
	return &common.Request{
		Wait: t.Wait,
		Body: t.Messages[t.index],
	}
}
