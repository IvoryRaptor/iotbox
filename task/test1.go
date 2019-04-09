package task

import (
	"github.com/IvoryRaptor/iotbox/akka"
	"time"
)

type PingActor struct {
	akka.Actor
}

func (actor *PingActor) Receive(sender akka.IActor, message akka.Message) error {
	println("PingActor")
	println(&message)
	actor.ActorSelect("pang").Tell(actor, message)
	return nil
}

func (actor *PingActor) PreStart() error {
	return nil
}

type PangActor struct {
	akka.Actor
}

func (actor *PangActor) PreStart() error {
	return nil
}

func (actor *PangActor) Receive(sender akka.IActor, message akka.Message) {
	time.Sleep(time.Second)
	println("PangActor")
	println(&message)
	//sender.Tell(actor, message)
}
func test1(system *akka.System) {
	//pingActor := system.ActorOf(&PingActor{}, "ping")
	//system.ActorOf(&PangActor{}, "pang")
	//pingActor.Tell(pingActor, akka.Message{})
}
