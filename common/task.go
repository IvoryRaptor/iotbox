package common

import (
	"github.com/IvoryRaptor/iotbox/akka"
	"time"
)

type Message map[string]interface{}

type Request struct {
	Wait time.Duration
	Body Message
}

type ResponseState int

const (
	ResponseResult  ResponseState = iota // value --> 0
	ResponseTimeout                      // value --> 1
)

type Response struct {
	State ResponseState
	Body  Message
}

type ITask interface {
	Init(task *TaskRef)
	Receive(task *TaskRef, response *Response)
	GetRequest(task *TaskRef) *Request
}

type TaskRef struct {
	task            ITask
	data            map[string]interface{}
	func_receive    func(task *TaskRef, response *Response)
	func_getrequest func(task *TaskRef) *Request
}

func (t *TaskRef) Receive(response *Response) {
	t.func_receive(t, response)
}

func (t *TaskRef) Become(module string, receive func(task *TaskRef, response *Response), getrequest func(task *TaskRef) *Request) {
	t.func_receive = receive
	t.func_getrequest = getrequest
	t.JoinModule(module)
}

func (t *TaskRef) GetRequest() *Request {
	return t.func_getrequest(t)
}
func (t *TaskRef) JoinModule(name string) {
	m := akka.NewLocalActorOf(name)
	akka.EmptyRootContext.Tell(m, t)
}
