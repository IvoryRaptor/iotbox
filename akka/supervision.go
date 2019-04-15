package akka

import (
	"time"

	"github.com/AsynkronIT/protoactor-go/eventstream"
)

// DeciderFunc is a function which is called by a SupervisorStrategy
type DeciderFunc func(reason interface{}) Directive

// SupervisorStrategy is an interface that decides how to handle failing child actors
type SupervisorStrategy interface {
	HandleFailure(supervisor Supervisor, child *ActorRef, rs *RestartStatistics, reason interface{}, message interface{})
}

// Supervisor is an interface that is used by the SupervisorStrategy to manage child actor lifecycle
type Supervisor interface {
	Children() []*ActorRef
	EscalateFailure(reason interface{}, message interface{})
	RestartChildren(actorOfs ...*ActorRef)
	StopChildren(actorOfs ...*ActorRef)
	ResumeChildren(actorOfs ...*ActorRef)
}

func logFailure(child *ActorRef, reason interface{}, directive Directive) {
	eventstream.Publish(&SupervisorEvent{
		Child:     child,
		Reason:    reason,
		Directive: directive,
	})
}

// DefaultDecider is a decider that will always restart the failing child actor
func DefaultDecider(_ interface{}) Directive {
	return RestartDirective
}

var (
	defaultSupervisionStrategy    = NewOneForOneStrategy(10, 10*time.Second, DefaultDecider)
	restartingSupervisionStrategy = NewRestartingStrategy()
)

func DefaultSupervisorStrategy() SupervisorStrategy {
	return defaultSupervisionStrategy
}

func RestartingSupervisorStrategy() SupervisorStrategy {
	return restartingSupervisionStrategy
}
