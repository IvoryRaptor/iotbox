package akka

import (
	"time"
)

type IActor interface {
	start(self IActor, system *System, queue chan *block)
	Receive(sender IActor, message Message)
	Tell(owner IActor, message Message)
	ActorSelect(path string) IActor
	Ask(owner IActor, message Message, timeOut time.Duration) (bool, Message)
}

type Actor struct {
	system *System
	self   IActor
}

func (actor *Actor) Tell(owner IActor, message Message) {
	actor.system.tell(owner, actor.self, message)
}

func (actor *Actor) Ask(owner IActor, message Message, timeOut time.Duration) (bool, Message) {
	return actor.system.ask(owner, actor.self, message, timeOut)
}

func (actor *Actor) ActorSelect(path string) IActor {
	return actor.system.paths[path]
}

func (actor *Actor) start(self IActor, system *System, queue chan *block) {
	actor.system = system
	actor.self = self
	go func() {
		for {
			var block = <-queue
			self.Receive(block.owner, block.message)
		}
	}()
}
