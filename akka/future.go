package akka

import (
	"errors"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/log"
)

// ErrTimeout is the error used when a future times out before receiving a result.
var ErrTimeout = errors.New("future: timeout")

// NewFuture creates and returns a new actor.Future with a timeout of duration d
func NewFuture(d time.Duration) *Future {
	ref := &futureProcess{Future{cond: sync.NewCond(&sync.Mutex{})}}
	id := ProcessRegistry.NextId()

	actorOf, ok := ProcessRegistry.Add(ref, "future"+id)
	if !ok {
		plog.Error("failed to register future process", log.Stringer(" actorOf", actorOf))
	}

	ref.actorOf = actorOf
	if d >= 0 {
		ref.t = time.AfterFunc(d, func() {
			ref.err = ErrTimeout
			ref.Stop(actorOf)
		})
	}

	return &ref.Future
}

type Future struct {
	actorOf *ActorRef
	cond    *sync.Cond
	// protected by cond
	done        bool
	result      interface{}
	err         error
	t           *time.Timer
	pipes       []*ActorRef
	completions []func(res interface{}, err error)
}

// ActorRef to the backing actor for the Future result
func (f *Future) ActorOf() *ActorRef {
	return f.actorOf
}

// PipeTo forwards the result or error of the future to the specified  actorOfs
func (f *Future) PipeTo(actorOfs ...*ActorRef) {
	f.cond.L.Lock()
	f.pipes = append(f.pipes, actorOfs...)
	// for an already completed future, force push the result to targets
	if f.done {
		f.sendToPipes()
	}
	f.cond.L.Unlock()
}

func (f *Future) sendToPipes() {
	if f.pipes == nil {
		return
	}

	var m interface{}
	if f.err != nil {
		m = f.err
	} else {
		m = f.result
	}
	for _, actorOf := range f.pipes {
		actorOf.sendUserMessage(m)
	}
	f.pipes = nil
}

func (f *Future) wait() {
	f.cond.L.Lock()
	for !f.done {
		f.cond.Wait()
	}
	f.cond.L.Unlock()
}

// Result waits for the future to resolve
func (f *Future) Result() (interface{}, error) {
	f.wait()
	return f.result, f.err
}

func (f *Future) Wait() error {
	f.wait()
	return f.err
}

func (f *Future) continueWith(continuation func(res interface{}, err error)) {
	f.cond.L.Lock()
	defer f.cond.L.Unlock() // use defer as the continuation could blow up
	if f.done {
		continuation(f.result, f.err)
	} else {
		f.completions = append(f.completions, continuation)
	}
}

// futureProcess is a struct carrying a response ActorRef and a channel where the response is placed
type futureProcess struct {
	Future
}

func (ref *futureProcess) SendUserMessage(actorOf *ActorRef, message interface{}) {
	_, msg, _ := UnwrapEnvelope(message)
	ref.result = msg
	ref.Stop(actorOf)
}

func (ref *futureProcess) SendSystemMessage(actorOf *ActorRef, message interface{}) {
	ref.result = message
	ref.Stop(actorOf)
}

func (ref *futureProcess) Stop(actorOf *ActorRef) {
	ref.cond.L.Lock()
	if ref.done {
		ref.cond.L.Unlock()
		return
	}

	ref.done = true
	if ref.t != nil {
		ref.t.Stop()
	}
	ProcessRegistry.Remove(actorOf)

	ref.sendToPipes()
	ref.runCompletions()
	ref.cond.L.Unlock()
	ref.cond.Signal()
}

// TODO: we could replace "pipes" with this
// instead of pushing ActorOfs to pipes, we could push wrapper funcs that tells the  actorOf
// as a completion, that would unify the model
func (f *Future) runCompletions() {
	if f.completions == nil {
		return
	}

	for _, c := range f.completions {
		c(f.result, f.err)
	}
	f.completions = nil
}
