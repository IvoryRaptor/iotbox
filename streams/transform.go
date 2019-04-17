package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Transform struct {
	Flow
	work func(msg interface{}) interface{}
}

func (t *Transform) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		t.start(context)
	case interface{}:
		res := t.work(msg)
		for _, next := range t.refs {
			context.Tell(next, res)
		}
	}
}

func (t *Transform) Map(work func(msg interface{}) interface{}) *Transform {
	return &Transform{
		work: work,
	}
}
