package akka

import (
	"fmt"
	"log"
	"time"
)

type IActor interface {
	Receive(sender IActor, message Message) error
	Tell(owner IActor, message Message)
	ActorSelect(path string) IActor
	Ask(owner IActor, message Message, timeOut time.Duration) (bool, Message)
	PreStart() error
	GetPath() string
	init(actor IActor, system *System, name string)
	start(queue chan *block) error
}

type Actor struct {
	system  *System
	self    IActor
	path    string
	receive func(sender IActor, message Message) error
}

func (actor *Actor) Info(format string, v ...interface{}) {
	log.Printf("[%s]====> %s", actor.GetPath(), fmt.Sprintf(format, v))
}

func (actor *Actor) Error(err error) {
	log.Printf("[%s] Error %s", actor.GetPath(), err.Error())
}

func (actor *Actor) Warn(format string, v ...interface{}) {
	log.Printf("[%s] Warn %s", actor.GetPath(), fmt.Sprintf(format, v))
}

func (actor *Actor) GetPath() string {
	return actor.path
}
func (actor *Actor) Tell(owner IActor, message Message) {
	actor.system.tell(owner, actor.self, message)
}

func (actor *Actor) Ask(owner IActor, message Message, timeOut time.Duration) (bool, Message) {
	return actor.system.ask(owner, actor.self, message, timeOut)
}

func (actor *Actor) Become(newReceive func(sender IActor, message Message) error) {
	actor.receive = newReceive
}

func (actor *Actor) ActorSelect(path string) IActor {
	return actor.system.paths[path]
}

func (actor *Actor) init(self IActor, system *System, path string) {
	actor.system = system
	actor.self = self
	actor.path = path
	actor.receive = self.Receive
}

func (actor *Actor) start(queue chan *block) error {
	if err := actor.self.PreStart(); err != nil {
		return err
	}
	go func() {
		for {
			var block = <-queue
			if err := actor.receive(block.owner, block.message); err != nil {
				log.Printf("[%s] Error %s", actor.GetPath(), err.Error())
			}
		}
	}()
	return nil
}
