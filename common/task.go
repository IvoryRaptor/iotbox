package common

type Message map[string]interface{}

type ITask interface {
	Init() Message
	GetNext(msg Message) Message
}
