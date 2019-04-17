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

func (w *Window) Map(work func(msg []interface{}) interface{}) *Transform {
	result := &Transform{
		work: func(msg interface{}) interface{} {
			return work(msg.([]interface{}))
		},
	}
	w.actors = append(w.actors, result)
	return result
}

func (w *Window) Foreach(work func(msg []interface{})) *Sink {
	result := &Sink{
		work: func(msg interface{}) {
			work(msg.([]interface{}))
		},
	}
	w.actors = append(w.actors, result)
	return result
}

func (w *Window) Filter(work func(msg []interface{}) bool) *Filter {
	result := &Filter{
		work: func(msg interface{}) bool {
			return work(msg.([]interface{}))
		},
	}
	w.actors = append(w.actors, result)
	return result
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
