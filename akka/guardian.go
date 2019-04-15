package akka

import (
	"errors"
	"sync"

	"github.com/AsynkronIT/protoactor-go/log"
)

type guardiansValue struct {
	guardians *sync.Map
}

var guardians = &guardiansValue{&sync.Map{}}

func (gs *guardiansValue) getGuardianActorOf(s SupervisorStrategy) *ActorRef {
	if g, ok := gs.guardians.Load(s); ok {
		return g.(*guardianProcess).actorOf
	}
	g := gs.newGuardian(s)
	gs.guardians.Store(s, g)
	return g.actorOf
}

// newGuardian creates and returns a new actor.guardianProcess with a timeout of duration d
func (gs *guardiansValue) newGuardian(s SupervisorStrategy) *guardianProcess {
	ref := &guardianProcess{strategy: s}
	id := ProcessRegistry.NextId()

	actorOf, ok := ProcessRegistry.Add(ref, "guardian"+id)
	if !ok {
		plog.Error("failed to register guardian process", log.Stringer(" actorOf", actorOf))
	}

	ref.actorOf = actorOf
	return ref
}

type guardianProcess struct {
	actorOf  *ActorRef
	strategy SupervisorStrategy
}

func (g *guardianProcess) SendUserMessage(actorOf *ActorRef, message interface{}) {
	panic(errors.New("Guardian actor cannot receive any user messages"))
}

func (g *guardianProcess) SendSystemMessage(actorOf *ActorRef, message interface{}) {
	if msg, ok := message.(*Failure); ok {
		g.strategy.HandleFailure(g, msg.Who, msg.RestartStats, msg.Reason, msg.Message)
	}
}

func (g *guardianProcess) Stop(actorOf *ActorRef) {
	// Ignore
}

func (g *guardianProcess) Children() []*ActorRef {
	panic(errors.New("Guardian does not hold its children ActorOfs"))
}

func (*guardianProcess) EscalateFailure(reason interface{}, message interface{}) {
	panic(errors.New("Guardian cannot escalate failure"))
}

func (*guardianProcess) RestartChildren(actorOfs ...*ActorRef) {
	for _, actorOf := range actorOfs {
		actorOf.sendSystemMessage(restartMessage)
	}
}

func (*guardianProcess) StopChildren(actorOfs ...*ActorRef) {
	for _, actorOf := range actorOfs {
		actorOf.sendSystemMessage(stopMessage)
	}
}

func (*guardianProcess) ResumeChildren(actorOfs ...*ActorRef) {
	for _, actorOf := range actorOfs {
		actorOf.sendSystemMessage(resumeMailboxMessage)
	}
}
