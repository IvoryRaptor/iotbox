package akka

import "time"

// Context contains contextual information for actors
type Context interface {
	infoPart
	basePart
	senderPart
	receiverPart
	actorOfPart
}

type SenderContext interface {
	infoPart
	senderPart
}

type ReceiverContext interface {
	infoPart
	receiverPart
}

type infoPart interface {
	// Parent returns the ActorRef for the current actors parent
	Parent() *ActorRef

	// Self returns the ActorRef for the current actor
	Self() *ActorRef

	// Actor returns the actor associated with this context
	Actor() Actor
}

type basePart interface {
	// ReceiveTimeout returns the current timeout
	ReceiveTimeout() time.Duration

	// Returns a slice of the actors children
	Children() []*ActorRef

	// Respond sends a response to the to the current `Sender`
	// If the Sender is nil, the actor will panic
	Respond(response interface{})

	// Stash stashes the current message on a stack for reprocessing when the actor restarts
	Stash()

	// Watch registers the actor as a monitor for the specified ActorRef
	Watch(actorOf *ActorRef)

	// Unwatch unregisters the actor as a monitor for the specified ActorRef
	Unwatch(actorOf *ActorRef)

	// SetReceiveTimeout sets the inactivity timeout, after which a ReceiveTimeout message will be sent to the actor.
	// A duration of less than 1ms will disable the inactivity timer.
	//
	// If a message is received before the duration d, the timer will be reset. If the message conforms to
	// the NotInfluenceReceiveTimeout interface, the timer will not be reset
	SetReceiveTimeout(d time.Duration)

	CancelReceiveTimeout()

	// Forward forwards current message to the given ActorRef
	Forward(actorOf *ActorRef)

	AwaitFuture(f *Future, continuation func(res interface{}, err error))
}

type senderPart interface {
	// Sender returns the ActorRef of actor that sent currently processed message
	Sender() *ActorRef

	// Message returns the current message to be processed
	Message() interface{}

	// MessageHeader returns the meta information for the currently processed message
	MessageHeader() ReadonlyMessageHeader

	// Tell sends a message to the given ActorRef
	Tell(actorOf *ActorRef, message interface{})

	// Request sends a message to the given ActorRef
	Request(actorOf *ActorRef, message interface{})

	// Request sends a message to the given ActorRef and also provides a Sender ActorRef
	RequestWithCustomSender(actorOf *ActorRef, message interface{}, sender *ActorRef)

	// Ask sends a message to a given ActorRef and returns a Future
	Ask(actorOf *ActorRef, message interface{}, timeout time.Duration) *Future
}

type receiverPart interface {
	Receive(envelope *MessageEnvelope)
}

type actorOfPart interface {
	// ActorOf starts a new child actor based on props and named with a unique id
	ActorOf(props *Props) *ActorRef

	// ActorOfPrefix starts a new child actor based on props and named using a prefix followed by a unique id
	ActorOfPrefix(props *Props, prefix string) *ActorRef

	// ActorOfNamed starts a new child actor based on props and named using the specified name
	//
	// ErrNameExists will be returned if id already exists
	//
	// Please do not use name sharing same pattern with system actors, for example "YourPrefix$1", "Remote$1", "future$1"
	ActorOfNamed(props *Props, id string) (*ActorRef, error)
}
