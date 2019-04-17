package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Window struct {
	Flow
	count  int
	index  int
	packet []interface{}
}

func (w *Window) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		w.index = 0
		w.packet = make([]interface{}, w.count)
		w.start(context)
	case interface{}:
		w.packet[w.index] = msg
		w.index++
		if w.index >= w.count {
			w.index = 0
			packet := make([]interface{}, w.count)
			copy(packet, w.packet)
			w.tellNext(packet)
		}
	}
}
