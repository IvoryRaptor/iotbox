package main

import (
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type ReadTask struct {
	Owner   string
	Message common.Message
	result  *common.Response
	Wait    time.Duration
	send    bool
}

func (t *ReadTask) Init(task *common.TaskRef) {
	t.send = false
}

func (t *ReadTask) OwnerReceive(task *common.TaskRef, response *common.Response) {
	t.send = true
	println("result")
}

func (t *ReadTask) OwnerGetRequest(task *common.TaskRef) *common.Request {
	if t.send {
		return nil
	}
	println("OwnerGetRequest")
	return &common.Request{
		Wait:  t.Wait,
		Body:  t.result.Body,
		Retry: 3,
	}
}

func (t *ReadTask) Receive(task *common.TaskRef, response *common.Response) {
	t.result = response
	task.Become(t.Owner, t.OwnerReceive, t.OwnerGetRequest)
}

func (t *ReadTask) GetRequest(task *common.TaskRef) *common.Request {
	return &common.Request{
		Wait:  t.Wait,
		Body:  t.Message,
		Retry: 3,
	}
}
