package streams

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Flow struct {
	self   akka.ActorRef
	refs   []*akka.ActorRef
	actors []akka.Actor
}

func (f *Flow) Start(context akka.Context) {
	f.refs = make([]*akka.ActorRef, len(f.actors))
	for i, next := range f.actors {
		f.refs[i] = context.ActorOf(akka.PropsFromProducer(func() akka.Actor {
			return next
		}))
	}
}

func (f *Flow) TellNext(msg interface{}) {
	for _, next := range f.refs {
		akka.EmptyRootContext.Tell(next, msg)
	}
}
func (f *Flow) Append(actor akka.Actor) {
	f.actors = append(f.actors, actor)
}
