package akka

import (
	//	"fmt"
	//	"github.com/gogo/protobuf/jsonpb"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

type ActorRef struct {
	Address string `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	Id      string `protobuf:"bytes,2,opt,name=Id,proto3" json:"Id,omitempty"`

	p *Process
}

/*
func (m *ActorRef) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	str := fmt.Sprintf("{\"Address\":\"%v\", \"Id\":\"%v\"}", m.Address, m.Id)
	return []byte(str), nil
}*/

func (actorOf *ActorRef) ref() Process {
	p := (*Process)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&actorOf.p))))
	if p != nil {
		if l, ok := (*p).(*ActorProcess); ok && atomic.LoadInt32(&l.dead) == 1 {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&actorOf.p)), nil)
		} else {
			return *p
		}
	}

	ref, exists := ProcessRegistry.Get(actorOf)
	if exists {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&actorOf.p)), unsafe.Pointer(&ref))
	}

	return ref
}

// sendUserMessage sends a messages asynchronously to the ActorRef
func (actorOf *ActorRef) sendUserMessage(message interface{}) {
	actorOf.ref().SendUserMessage(actorOf, message)
}

func (actorOf *ActorRef) sendSystemMessage(message interface{}) {
	actorOf.ref().SendSystemMessage(actorOf, message)
}

// StopFuture will stop actor immediately regardless of existing user messages in mailbox, and return its future.
func (actorOf *ActorRef) StopFuture() *Future {
	future := NewFuture(10 * time.Second)

	actorOf.sendSystemMessage(&Watch{Watcher: future.actorOf})
	actorOf.Stop()

	return future
}

// GracefulStop will wait actor to stop immediately regardless of existing user messages in mailbox
func (actorOf *ActorRef) GracefulStop() {
	actorOf.StopFuture().Wait()
}

// Stop will stop actor immediately regardless of existing user messages in mailbox.
func (actorOf *ActorRef) Stop() {
	actorOf.ref().Stop(actorOf)
}

// PoisonFuture will tell actor to stop after processing current user messages in mailbox, and return its future.
func (actorOf *ActorRef) PoisonFuture() *Future {
	future := NewFuture(10 * time.Second)

	actorOf.sendSystemMessage(&Watch{Watcher: future.actorOf})
	actorOf.Poison()

	return future
}

// GracefulPoison will tell and wait actor to stop after processing current user messages in mailbox.
func (actorOf *ActorRef) GracefulPoison() {
	actorOf.PoisonFuture().Wait()
}

// Poison will tell actor to stop after processing current user messages in mailbox.
func (actorOf *ActorRef) Poison() {
	actorOf.sendUserMessage(&PoisonPill{})
}

func (actorOf *ActorRef) key() string {
	if actorOf.Address == ProcessRegistry.Address {
		return actorOf.Id
	}
	return actorOf.Address + "#" + actorOf.Id
}

func (actorOf *ActorRef) String() string {
	if actorOf == nil {
		return "nil"
	}
	return actorOf.Address + "/" + actorOf.Id
}

// NewActorOf returns a new instance of the ActorRef struct
func NewActorOf(address, id string) *ActorRef {
	return &ActorRef{
		Address: address,
		Id:      id,
	}
}

// NewLocalActorOf returns a new instance of the ActorRef struct with the address preset
func NewLocalActorOf(id string) *ActorRef {
	return &ActorRef{
		Address: ProcessRegistry.Address,
		Id:      id,
	}
}

func actorOfFromKey(key string, p *ActorRef) {
	i := strings.IndexByte(key, '#')
	if i == -1 {
		p.Address = ProcessRegistry.Address
		p.Id = key
	} else {
		p.Address = key[:i]
		p.Id = key[i+1:]
	}
}

// Deprecated: Use Context.Tell instead
func (actorOf *ActorRef) Tell(message interface{}) {
	ctx := EmptyRootContext
	ctx.Tell(actorOf, message)
}

// Deprecated: Use Context.Request or Context.RequestWithCustomSender instead
func (actorOf *ActorRef) Request(message interface{}, respondTo *ActorRef) {
	ctx := EmptyRootContext
	ctx.RequestWithCustomSender(actorOf, message, respondTo)
}

// Deprecated: Use Context.Ask instead
func (actorOf *ActorRef) Ask(message interface{}, timeout time.Duration) *Future {
	ctx := EmptyRootContext
	return ctx.Ask(actorOf, message, timeout)
}
