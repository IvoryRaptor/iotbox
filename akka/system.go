package akka

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type System struct {
	Actor
	paths  map[string]IActor
	queues map[IActor]chan *block
}

func (system *System) PreStart() error {
	return nil
}

func (system *System) Start() {
	system.paths = map[string]IActor{}
	system.queues = map[IActor]chan *block{}
	system.ActorOf(system, "System")
}

func (system *System) Receive(actor IActor, message Message) error {
	return nil
}

func (system *System) tell(from IActor, to IActor, message Message) {
	system.queues[to] <- &block{
		owner:   from,
		message: message,
	}
}

func (system *System) ask(from IActor, to IActor, message Message, timeOut time.Duration) (bool, Message) {
	actor := &AskActor{}
	system.ActorOf(actor, "")
	defer system.Remove(actor)
	system.queues[to] <- &block{
		owner:   actor,
		message: message,
	}
	time.AfterFunc(timeOut, func() {
		system.queues[actor] <- &block{
			owner:   system,
			message: nil,
		}
	})
	packet := <-actor.result
	return packet != nil, packet
}

func (system *System) Remove(actor IActor) {
	delete(system.queues, actor)
}

func (system *System) config(c IConfig, name string) error {
	file := fmt.Sprintf("./config/%s.yaml", name)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	return c.Config(config)
}

func (system *System) ActorOf(actor IActor, name string) IActor {
	actor.init(actor, system, name)
	if config, b := actor.(IConfig); b {
		system.config(config, name)
	}
	if name != "" {
		system.paths[name] = actor
	}
	queue := make(chan *block, 10)
	system.queues[actor] = queue
	if err := actor.start(queue); err != nil {
		log.Printf("[%s]====> Exec %s", name, err.Error())
		return nil
	}
	return actor
}
