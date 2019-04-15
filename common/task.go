package common

import "time"

type Message map[string]interface{}

type Request struct {
	Wait time.Duration
	Body Message
}

type State int

const (
	Initialize State = iota // value --> 0
	Result                  // value --> 1
	Timeout                 // value --> 2
)

type Response struct {
	State State
	Body  Message
}

type ITask interface {
	GetNext(response *Response) *Request
}
