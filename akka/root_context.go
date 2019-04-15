package akka

import "time"

// RootContext a Context can be used outside of Actor
type RootContext struct {
	senderMiddleware SenderFunc
	headers          messageHeader
}

// EmptyRootContext returns the default RootContext.
// Please do not set any headers/middlewares to this context.
var EmptyRootContext = &RootContext{
	senderMiddleware: nil,
	headers:          EmptyMessageHeader,
}

// NewRootContext creates a new RootContext that can be customized with headers and middlewares.
func NewRootContext(header map[string]string, middleware ...SenderMiddleware) *RootContext {
	return &RootContext{
		senderMiddleware: makeSenderMiddlewareChain(middleware, func(_ SenderContext, target *ActorRef, envelope *MessageEnvelope) {
			target.sendUserMessage(envelope)
		}),
		headers: messageHeader(header),
	}
}

// WithHeaders set headers to RootContext.
func (rc *RootContext) WithHeaders(headers map[string]string) *RootContext {
	rc.headers = headers
	return rc
}

// WithSenderMiddleware set SenderMiddlewares to RootContext.
func (rc *RootContext) WithSenderMiddleware(middleware ...SenderMiddleware) *RootContext {
	rc.senderMiddleware = makeSenderMiddlewareChain(middleware, func(_ SenderContext, target *ActorRef, envelope *MessageEnvelope) {
		target.sendUserMessage(envelope)
	})
	return rc
}

//
// Interface: info
//

func (rc *RootContext) Parent() *ActorRef {
	return nil
}

func (rc *RootContext) Self() *ActorRef {
	return nil
}

func (rc *RootContext) Sender() *ActorRef {
	return nil
}

func (rc *RootContext) Actor() Actor {
	return nil
}

//
// Interface: sender
//

func (rc *RootContext) Message() interface{} {
	return nil
}

func (rc *RootContext) MessageHeader() ReadonlyMessageHeader {
	return rc.headers
}

func (rc *RootContext) Tell(actorOf *ActorRef, message interface{}) {
	rc.sendUserMessage(actorOf, message)
}

func (rc *RootContext) Request(actorOf *ActorRef, message interface{}) {
	rc.sendUserMessage(actorOf, message)
}

func (rc *RootContext) RequestWithCustomSender(actorOf *ActorRef, message interface{}, sender *ActorRef) {
	env := &MessageEnvelope{
		Header:  nil,
		Message: message,
		Sender:  sender,
	}
	rc.sendUserMessage(actorOf, env)
}

// Ask sends a message to a given ActorRef and returns a Future
func (rc *RootContext) Ask(actorOf *ActorRef, message interface{}, timeout time.Duration) *Future {
	future := NewFuture(timeout)
	env := &MessageEnvelope{
		Header:  nil,
		Message: message,
		Sender:  future.ActorOf(),
	}
	rc.sendUserMessage(actorOf, env)
	return future
}

func (rc *RootContext) sendUserMessage(actorOf *ActorRef, message interface{}) {
	if rc.senderMiddleware != nil {
		if envelope, ok := message.(*MessageEnvelope); ok {
			// Request based middleware
			rc.senderMiddleware(rc, actorOf, envelope)
		} else {
			// tell based middleware
			rc.senderMiddleware(rc, actorOf, &MessageEnvelope{nil, message, nil})
		}
		return
	}
	// Default path
	actorOf.sendUserMessage(message)
}

//
// Interface: actorOfer
//

// ActorOf starts a new actor based on props and named with a unique id
func (rc *RootContext) ActorOf(props *Props) *ActorRef {
	actorOf, err := rc.ActorOfNamed(props, ProcessRegistry.NextId())
	if err != nil {
		panic(err)
	}
	return actorOf
}

// ActorOfPrefix starts a new actor based on props and named using a prefix followed by a unique id
func (rc *RootContext) ActorOfPrefix(props *Props, prefix string) *ActorRef {
	actorOf, err := rc.ActorOfNamed(props, prefix+ProcessRegistry.NextId())
	if err != nil {
		panic(err)
	}
	return actorOf
}

// ActorOfNamed starts a new actor based on props and named using the specified name
//
// ErrNameExists will be returned if id already exists
//
// Please do not use name sharing same pattern with system actors, for example "YourPrefix$1", "Remote$1", "future$1"
func (rc *RootContext) ActorOfNamed(props *Props, name string) (*ActorRef, error) {
	var parent *ActorRef
	if props.guardianStrategy != nil {
		parent = guardians.getGuardianActorOf(props.guardianStrategy)
	}
	return props.actorOf(name, parent)
}
