package common

import "github.com/IvoryRaptor/iotbox/akka"

type Protocol interface {
}

type Port struct {
	protocol Protocol
}

func (port *Port) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case Message:
		println(msg["name"].(string))
		response := Message{
			"name":  msg["name"].(string),
			"value": "1",
		}
		context.Tell(context.Sender(), response)
	}
}
