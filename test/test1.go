package main

import (
	"github.com/IvoryRaptor/iotbox/test/akka"
	"time"
)

type PingActor struct {
	akka.Actor
}

func (actor *PingActor) Receive(sender akka.IActor, message akka.Message) {
	println("PingActor")
	println(&message)
	actor.ActorSelect("pang").Tell(actor, message)
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
