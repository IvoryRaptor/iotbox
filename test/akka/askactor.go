package akka

type AskActor struct {
	Actor
	result chan Message
}

func (actor *AskActor) Receive(sender IActor, message Message) {
	actor.result <- message
}

func (actor *AskActor) PreStart() error {
	actor.result = make(chan Message, 1)
	return nil
}
