package main

import (
	"github.com/IvoryRaptor/iotbox/test/akka"
	"time"
)

type ServerActor struct {
	akka.Actor
}

func (actor *ServerActor) Receive(sender akka.IActor, message akka.Message) {
	msg := akka.Message{}
	msg["result"] = message.GetInt("a") + message.GetInt("b")
	//time.Sleep(time.Second * 3)
	sender.Tell(actor, msg)
}

type ClientActor struct {
	akka.Actor
}

func (actor *ClientActor) Receive(_ akka.IActor, tmp akka.Message) {
	message := akka.Message{}
	message["a"] = 1
	message["b"] = 2
	finish, result := actor.ActorSelect("server").Ask(actor, message, time.Second*2)
	if finish {
		println(result.GetInt("result"))
	} else {
		println("timeout")
	}
}
func test2(system *akka.System) {
	system.ActorOf(&ServerActor{}, "server")
	client := system.ActorOf(&ClientActor{}, "client")
	client.Tell(nil, nil)
}
