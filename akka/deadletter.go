package akka

import (
	"github.com/AsynkronIT/protoactor-go/eventstream"
	"github.com/AsynkronIT/protoactor-go/log"
)

type deadLetterProcess struct{}

var (
	deadLetter           Process = &deadLetterProcess{}
	deadLetterSubscriber *eventstream.Subscription
)

func init() {
	deadLetterSubscriber = eventstream.Subscribe(func(msg interface{}) {
		if deadLetter, ok := msg.(*DeadLetterEvent); ok {
			plog.Debug("[DeadLetter]", log.Stringer(" actorOf", deadLetter.ActorOf), log.Message(deadLetter.Message), log.Stringer("sender", deadLetter.Sender))
		}
	})

	// this subscriber may not be deactivated.
	// it ensures that Watch commands that reach a stopped actor gets a Terminated message back.
	// This can happen if one actor tries to Watch a ActorRef, while another thread sends a Stop message.
	eventstream.Subscribe(func(msg interface{}) {
		if deadLetter, ok := msg.(*DeadLetterEvent); ok {
			if m, ok := deadLetter.Message.(*Watch); ok {
				// we know that this is a local actor since we get it on our own event stream, thus the address is not terminated
				m.Watcher.sendSystemMessage(&Terminated{AddressTerminated: false, Who: deadLetter.ActorOf})
			}
		}
	})
}

// A DeadLetterEvent is published via event.Publish when a message is sent to a nonexistent ActorRef
type DeadLetterEvent struct {
	ActorOf *ActorRef   // The invalid process, to which the message was sent
	Message interface{} // The message that could not be delivered
	Sender  *ActorRef   // the process that sent the Message
}

func (*deadLetterProcess) SendUserMessage(actorOf *ActorRef, message interface{}) {
	_, msg, sender := UnwrapEnvelope(message)
	eventstream.Publish(&DeadLetterEvent{
		ActorOf: actorOf,
		Message: msg,
		Sender:  sender,
	})
}

func (*deadLetterProcess) SendSystemMessage(actorOf *ActorRef, message interface{}) {
	eventstream.Publish(&DeadLetterEvent{
		ActorOf: actorOf,
		Message: message,
	})
}

func (ref *deadLetterProcess) Stop(actorOf *ActorRef) {
	ref.SendSystemMessage(actorOf, stopMessage)
}
