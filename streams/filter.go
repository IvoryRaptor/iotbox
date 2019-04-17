package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Filter struct {
	Flow
	work func(msg interface{}) bool
}

func (f *Filter) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		f.start(context)
	case interface{}:
		if f.work(msg) {
			f.tellNext(msg)
		}
	}
}
