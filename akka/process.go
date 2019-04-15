package akka

// A Process is an interface that defines the base contract for interaction of actors
type Process interface {
	SendUserMessage(actorOf *ActorRef, message interface{})
	SendSystemMessage(actorOf *ActorRef, message interface{})
	Stop(actorOf *ActorRef)
}
