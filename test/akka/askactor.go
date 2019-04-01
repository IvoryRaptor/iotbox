package akka

type AskActor struct {
	Actor
	result chan Message
}

func (actor *AskActor) Receive(sender IActor, message Message) {
	actor.result <- message
}
