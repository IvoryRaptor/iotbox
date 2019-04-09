package akka

type AskActor struct {
	Actor
	result chan Message
}

func (actor *AskActor) Receive(sender IActor, message Message) error {
	actor.result <- message
	return nil
}

func (actor *AskActor) PreStart() error {
	actor.result = make(chan Message, 1)
	return nil
}
