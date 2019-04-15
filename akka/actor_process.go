package akka

import (
	"sync/atomic"

	"github.com/AsynkronIT/protoactor-go/mailbox"
)

type ActorProcess struct {
	mailbox mailbox.Mailbox
	dead    int32
}

func NewActorProcess(mailbox mailbox.Mailbox) *ActorProcess {
	return &ActorProcess{mailbox: mailbox}
}

func (ref *ActorProcess) SendUserMessage(actorOf *ActorRef, message interface{}) {
	ref.mailbox.PostUserMessage(message)
}
func (ref *ActorProcess) SendSystemMessage(actorOf *ActorRef, message interface{}) {
	ref.mailbox.PostSystemMessage(message)
}

func (ref *ActorProcess) Stop(actorOf *ActorRef) {
	atomic.StoreInt32(&ref.dead, 1)
	ref.SendSystemMessage(actorOf, stopMessage)
}
