package common

import "time"

type Message map[string]interface{}

type Request struct {
	Wait time.Duration
	Msg  Message
}

type ITask interface {
	Init() *Request
	GetNext(msg Message) *Request
}
