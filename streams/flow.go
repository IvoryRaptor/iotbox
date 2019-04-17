package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Flow struct {
	refs   []*akka.ActorRef
	actors []akka.Actor
}

func (f *Flow) Map(work func(msg interface{}) interface{}) *Transform {
	result := &Transform{
		work: work,
	}
	f.actors = append(f.actors, result)
	return result
}

func (f *Flow) Foreach(work func(msg interface{})) *Sink {
	result := &Sink{
		work: work,
	}
	f.actors = append(f.actors, result)
	return result
}

func (f *Flow) Filter(work func(msg interface{}) bool) *Filter {
	result := &Filter{
		work: work,
	}
	f.actors = append(f.actors, result)
	return result
}

func (f *Flow) Window(count int) *Window {
	result := &Window{
		count: count,
	}
	f.actors = append(f.actors, result)
	return result
}

func (f *Flow) start(context akka.Context) {
	f.refs = make([]*akka.ActorRef, len(f.actors))
	for i, next := range f.actors {
		f.refs[i] = context.ActorOf(akka.PropsFromProducer(func() akka.Actor {
			return next
		}))
	}
}

func (f *Flow) tellNext(msg interface{}) {
	for _, next := range f.refs {
		akka.EmptyRootContext.Tell(next, msg)
	}
}
