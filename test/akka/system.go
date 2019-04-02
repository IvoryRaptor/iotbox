package akka

import "time"

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

func (system *System) Receive(actor IActor, message Message) {

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

func (system *System) ActorOf(actor IActor, name string) IActor {
	if name != "" {
		system.paths[name] = actor
	}
	queue := make(chan *block, 10)
	system.queues[actor] = queue
	actor.start(actor, system, queue)
	return actor
}
