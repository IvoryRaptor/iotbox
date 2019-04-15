package akka

import (
	"errors"
	"time"

	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/emirpasic/gods/stacks/linkedliststack"
)

type contextState int32

const (
	stateNone contextState = iota
	stateAlive
	stateRestarting
	stateStopping
	stateStopped
)

type actorContextExtras struct {
	children            ActorOfSet
	receiveTimeoutTimer *time.Timer
	rs                  *RestartStatistics
	stash               *linkedliststack.Stack
	watchers            ActorOfSet
	context             Context
}

func newActorContextExtras(context Context) *actorContextExtras {
	this := &actorContextExtras{
		context: context,
	}
	return this
}

func (ctxExt *actorContextExtras) restartStats() *RestartStatistics {
	// lazy initialize the child restart stats if this is the first time
	// further mutations are handled within "restart"
	if ctxExt.rs == nil {
		ctxExt.rs = NewRestartStatistics()
	}
	return ctxExt.rs
}

func (ctxExt *actorContextExtras) initReceiveTimeoutTimer(timer *time.Timer) {
	ctxExt.receiveTimeoutTimer = timer
}

func (ctxExt *actorContextExtras) resetReceiveTimeoutTimer(time time.Duration) {
	if ctxExt.receiveTimeoutTimer == nil {
		return
	}
	ctxExt.receiveTimeoutTimer.Reset(time)
}

func (ctxExt *actorContextExtras) stopReceiveTimeoutTimer() {
	if ctxExt.receiveTimeoutTimer == nil {
		return
	}
	ctxExt.receiveTimeoutTimer.Stop()
}

func (ctxExt *actorContextExtras) killReceiveTimeoutTimer() {
	if ctxExt.receiveTimeoutTimer == nil {
		return
	}
	ctxExt.receiveTimeoutTimer.Stop()
	ctxExt.receiveTimeoutTimer = nil
}

func (ctxExt *actorContextExtras) addChild(actorOf *ActorRef) {
	ctxExt.children.Add(actorOf)
}

func (ctxExt *actorContextExtras) removeChild(actorOf *ActorRef) {
	ctxExt.children.Remove(actorOf)
}

func (ctxExt *actorContextExtras) watch(watcher *ActorRef) {
	ctxExt.watchers.Add(watcher)
}

func (ctxExt *actorContextExtras) unwatch(watcher *ActorRef) {
	ctxExt.watchers.Remove(watcher)
}

type actorContext struct {
	actor             Actor
	extras            *actorContextExtras
	props             *Props
	parent            *ActorRef
	self              *ActorRef
	receiveTimeout    time.Duration
	producer          Producer
	messageOrEnvelope interface{}
	state             contextState
}

func newActorContext(props *Props, parent *ActorRef) *actorContext {
	this := &actorContext{
		parent: parent,
		props:  props,
	}

	this.incarnateActor()
	return this
}

func (ctx *actorContext) ensureExtras() *actorContextExtras {
	if ctx.extras == nil {
		ctxd := Context(ctx)
		if ctx.props != nil && ctx.props.contextDecoratorChain != nil {
			ctxd = ctx.props.contextDecoratorChain(ctxd)
		}
		ctx.extras = newActorContextExtras(ctxd)
	}
	return ctx.extras
}

//
// Interface: Context
//

func (ctx *actorContext) Parent() *ActorRef {
	return ctx.parent
}

func (ctx *actorContext) Self() *ActorRef {
	return ctx.self
}

func (ctx *actorContext) Sender() *ActorRef {
	return UnwrapEnvelopeSender(ctx.messageOrEnvelope)
}

func (ctx *actorContext) Actor() Actor {
	return ctx.actor
}

func (ctx *actorContext) ReceiveTimeout() time.Duration {
	return ctx.receiveTimeout
}

func (ctx *actorContext) Children() []*ActorRef {
	if ctx.extras == nil {
		return make([]*ActorRef, 0)
	}

	r := make([]*ActorRef, ctx.extras.children.Len())
	ctx.extras.children.ForEach(func(i int, p ActorRef) {
		r[i] = &p
	})
	return r
}

func (ctx *actorContext) Respond(response interface{}) {
	// If the message is addressed to nil forward it to the dead letter channel
	if ctx.Sender() == nil {
		deadLetter.SendUserMessage(nil, response)
		return
	}

	ctx.Tell(ctx.Sender(), response)
}

func (ctx *actorContext) Stash() {
	extra := ctx.ensureExtras()
	if extra.stash == nil {
		extra.stash = linkedliststack.New()
	}
	extra.stash.Push(ctx.Message())
}

func (ctx *actorContext) Watch(who *ActorRef) {
	who.sendSystemMessage(&Watch{
		Watcher: ctx.self,
	})
}

func (ctx *actorContext) Unwatch(who *ActorRef) {
	who.sendSystemMessage(&Unwatch{
		Watcher: ctx.self,
	})
}

func (ctx *actorContext) SetReceiveTimeout(d time.Duration) {
	if d <= 0 {
		panic("Duration must be greater than zero")
	}

	if d == ctx.receiveTimeout {
		return
	}

	if d < time.Millisecond {
		// anything less than than 1 millisecond is set to zero
		d = 0
	}

	ctx.receiveTimeout = d

	ctx.ensureExtras()
	ctx.extras.stopReceiveTimeoutTimer()
	if d > 0 {
		if ctx.extras.receiveTimeoutTimer == nil {
			ctx.extras.initReceiveTimeoutTimer(time.AfterFunc(d, ctx.receiveTimeoutHandler))
		} else {
			ctx.extras.resetReceiveTimeoutTimer(d)
		}
	}
}

func (ctx *actorContext) CancelReceiveTimeout() {
	if ctx.extras == nil || ctx.extras.receiveTimeoutTimer == nil {
		return
	}

	ctx.extras.killReceiveTimeoutTimer()
	ctx.receiveTimeout = 0
}

func (ctx *actorContext) receiveTimeoutHandler() {
	if ctx.extras != nil && ctx.extras.receiveTimeoutTimer != nil {
		ctx.CancelReceiveTimeout()
		ctx.Tell(ctx.self, receiveTimeoutMessage)
	}
}

func (ctx *actorContext) Forward(actorOf *ActorRef) {
	if msg, ok := ctx.messageOrEnvelope.(SystemMessage); ok {
		// SystemMessage cannot be forwarded
		plog.Error("SystemMessage cannot be forwarded", log.Message(msg))
		return
	}
	ctx.sendUserMessage(actorOf, ctx.messageOrEnvelope)
}

func (ctx *actorContext) AwaitFuture(f *Future, cont func(res interface{}, err error)) {
	wrapper := func() {
		cont(f.result, f.err)
	}

	message := ctx.messageOrEnvelope
	// invoke the callback when the future completes
	f.continueWith(func(res interface{}, err error) {
		// send the wrapped callaback as a continuation message to self
		ctx.self.sendSystemMessage(&continuation{
			f:       wrapper,
			message: message,
		})
	})
}

//
// Interface: sender
//

func (ctx *actorContext) Message() interface{} {
	return UnwrapEnvelopeMessage(ctx.messageOrEnvelope)
}

func (ctx *actorContext) MessageHeader() ReadonlyMessageHeader {
	return UnwrapEnvelopeHeader(ctx.messageOrEnvelope)
}

func (ctx *actorContext) Tell(actorOf *ActorRef, message interface{}) {
	ctx.sendUserMessage(actorOf, message)
}

func (ctx *actorContext) sendUserMessage(actorOf *ActorRef, message interface{}) {
	if ctx.props.senderMiddlewareChain != nil {
		ctx.props.senderMiddlewareChain(ctx.ensureExtras().context, actorOf, WrapEnvelope(message))
	} else {
		actorOf.sendUserMessage(message)
	}
}

func (ctx *actorContext) Request(actorOf *ActorRef, message interface{}) {
	env := &MessageEnvelope{
		Header:  nil,
		Message: message,
		Sender:  ctx.Self(),
	}

	ctx.sendUserMessage(actorOf, env)
}

func (rc *actorContext) RequestWithCustomSender(actorOf *ActorRef, message interface{}, sender *ActorRef) {
	env := &MessageEnvelope{
		Header:  nil,
		Message: message,
		Sender:  sender,
	}
	rc.sendUserMessage(actorOf, env)
}

func (ctx *actorContext) Ask(actorOf *ActorRef, message interface{}, timeout time.Duration) *Future {
	future := NewFuture(timeout)
	env := &MessageEnvelope{
		Header:  nil,
		Message: message,
		Sender:  future.ActorOf(),
	}
	ctx.sendUserMessage(actorOf, env)

	return future
}

//
// Interface: receiver
//

func (ctx *actorContext) Receive(envelope *MessageEnvelope) {
	ctx.messageOrEnvelope = envelope
	ctx.defaultReceive()
	ctx.messageOrEnvelope = nil
}

func (ctx *actorContext) defaultReceive() {
	if _, ok := ctx.Message().(*PoisonPill); ok {
		ctx.self.Stop()
		return
	}

	// are we using decorators, if so, ensure it has been created
	if ctx.props.contextDecoratorChain != nil {
		ctx.actor.Receive(ctx.ensureExtras().context)
		return
	}

	ctx.actor.Receive(Context(ctx))
}

//
// Interface: actorOfer
//

func (ctx *actorContext) ActorOf(props *Props) *ActorRef {
	actorOf, err := ctx.ActorOfNamed(props, ProcessRegistry.NextId())
	if err != nil {
		panic(err)
	}
	return actorOf
}

func (ctx *actorContext) ActorOfPrefix(props *Props, prefix string) *ActorRef {
	actorOf, err := ctx.ActorOfNamed(props, prefix+ProcessRegistry.NextId())
	if err != nil {
		panic(err)
	}
	return actorOf
}

func (ctx *actorContext) ActorOfNamed(props *Props, name string) (*ActorRef, error) {
	if props.guardianStrategy != nil {
		panic(errors.New("Props used to actorOf child cannot have GuardianStrategy"))
	}

	actorOf, err := props.actorOf(ctx.self.Id+"/"+name, ctx.self)
	if err != nil {
		return actorOf, err
	}

	ctx.ensureExtras().addChild(actorOf)

	return actorOf, nil
}

//
// Interface: MessageInvoker
//

func (ctx *actorContext) InvokeUserMessage(md interface{}) {
	if ctx.state == stateStopped {
		// already stopped
		return
	}

	influenceTimeout := true
	if ctx.receiveTimeout > 0 {
		_, influenceTimeout = md.(NotInfluenceReceiveTimeout)
		influenceTimeout = !influenceTimeout
		if influenceTimeout {
			ctx.extras.stopReceiveTimeoutTimer()
		}
	}

	ctx.processMessage(md)

	if ctx.receiveTimeout > 0 && influenceTimeout {
		ctx.extras.resetReceiveTimeoutTimer(ctx.receiveTimeout)
	}
}

func (ctx *actorContext) processMessage(m interface{}) {
	if ctx.props.receiverMiddlewareChain != nil {
		ctx.props.receiverMiddlewareChain(ctx.ensureExtras().context, WrapEnvelope(m))
		return
	}

	if ctx.props.contextDecoratorChain != nil {
		ctx.ensureExtras().context.Receive(WrapEnvelope(m))
		return
	}

	ctx.messageOrEnvelope = m
	ctx.defaultReceive()
	ctx.messageOrEnvelope = nil // release message
}

func (ctx *actorContext) incarnateActor() {
	ctx.state = stateAlive
	ctx.actor = ctx.props.producer()
}

func (ctx *actorContext) InvokeSystemMessage(message interface{}) {
	switch msg := message.(type) {
	case *continuation:
		ctx.messageOrEnvelope = msg.message // apply the message that was present when we started the await
		msg.f()                             // invoke the continuation in the current actor context
		ctx.messageOrEnvelope = nil         // release the message
	case *Started:
		ctx.InvokeUserMessage(msg) // forward
	case *Watch:
		ctx.handleWatch(msg)
	case *Unwatch:
		ctx.handleUnwatch(msg)
	case *Stop:
		ctx.handleStop(msg)
	case *Terminated:
		ctx.handleTerminated(msg)
	case *Failure:
		ctx.handleFailure(msg)
	case *Restart:
		ctx.handleRestart(msg)
	default:
		plog.Error("unknown system message", log.Message(msg))
	}
}

func (ctx *actorContext) handleRootFailure(failure *Failure) {
	defaultSupervisionStrategy.HandleFailure(ctx, failure.Who, failure.RestartStats, failure.Reason, failure.Message)
}

func (ctx *actorContext) handleWatch(msg *Watch) {
	if ctx.state >= stateStopping {
		msg.Watcher.sendSystemMessage(&Terminated{
			Who: ctx.self,
		})
	} else {
		ctx.ensureExtras().watch(msg.Watcher)
	}
}

func (ctx *actorContext) handleUnwatch(msg *Unwatch) {
	if ctx.extras == nil {
		return
	}
	ctx.extras.unwatch(msg.Watcher)
}

func (ctx *actorContext) handleRestart(msg *Restart) {
	ctx.state = stateRestarting
	ctx.InvokeUserMessage(restartingMessage)
	ctx.stopAllChildren()
	ctx.tryRestartOrTerminate()
}

// I am stopping
func (ctx *actorContext) handleStop(msg *Stop) {
	if ctx.state >= stateStopping {
		// already stopping or stopped
		return
	}

	ctx.state = stateStopping

	ctx.InvokeUserMessage(stoppingMessage)
	ctx.stopAllChildren()
	ctx.tryRestartOrTerminate()
}

// child stopped, check if we can stop or restart (if needed)
func (ctx *actorContext) handleTerminated(msg *Terminated) {
	if ctx.extras != nil {
		ctx.extras.removeChild(msg.Who)
	}

	ctx.InvokeUserMessage(msg)
	ctx.tryRestartOrTerminate()
}

// offload the supervision completely to the supervisor strategy
func (ctx *actorContext) handleFailure(msg *Failure) {
	if strategy, ok := ctx.actor.(SupervisorStrategy); ok {
		strategy.HandleFailure(ctx, msg.Who, msg.RestartStats, msg.Reason, msg.Message)
		return
	}
	ctx.props.getSupervisor().HandleFailure(ctx, msg.Who, msg.RestartStats, msg.Reason, msg.Message)
}

func (ctx *actorContext) stopAllChildren() {
	if ctx.extras == nil {
		return
	}
	ctx.extras.children.ForEach(func(_ int, actorOf ActorRef) {
		actorOf.Stop()
	})
}

func (ctx *actorContext) tryRestartOrTerminate() {
	if ctx.extras != nil && !ctx.extras.children.Empty() {
		return
	}

	ctx.CancelReceiveTimeout()

	switch ctx.state {
	case stateRestarting:
		ctx.restart()
	case stateStopping:
		ctx.finalizeStop()
	}
}

func (ctx *actorContext) restart() {
	ctx.incarnateActor()
	ctx.self.sendSystemMessage(resumeMailboxMessage)
	ctx.InvokeUserMessage(startedMessage)
	if ctx.extras != nil && ctx.extras.stash != nil {
		for !ctx.extras.stash.Empty() {
			msg, _ := ctx.extras.stash.Pop()
			ctx.InvokeUserMessage(msg)
		}
	}
}

func (ctx *actorContext) finalizeStop() {
	ProcessRegistry.Remove(ctx.self)
	ctx.InvokeUserMessage(stoppedMessage)
	otherStopped := &Terminated{Who: ctx.self}
	// Notify watchers
	if ctx.extras != nil {
		ctx.extras.watchers.ForEach(func(i int, actorOf ActorRef) {
			actorOf.sendSystemMessage(otherStopped)
		})
	}
	// Notify parent
	if ctx.parent != nil {
		ctx.parent.sendSystemMessage(otherStopped)
	}
	ctx.state = stateStopped
}

//
// Interface: Supervisor
//

func (ctx *actorContext) EscalateFailure(reason interface{}, message interface{}) {
	failure := &Failure{Reason: reason, Who: ctx.self, RestartStats: ctx.ensureExtras().restartStats(), Message: message}
	ctx.self.sendSystemMessage(suspendMailboxMessage)
	if ctx.parent == nil {
		ctx.handleRootFailure(failure)
	} else {
		// TODO: Akka recursively suspends all children also on failure
		// Not sure if I think this is the right way to go, why do children need to wait for their parents failed state to recover?
		ctx.parent.sendSystemMessage(failure)
	}
}

func (*actorContext) RestartChildren(actorOfs ...*ActorRef) {
	for _, actorOf := range actorOfs {
		actorOf.sendSystemMessage(restartMessage)
	}
}

func (*actorContext) StopChildren(actorOfs ...*ActorRef) {
	for _, actorOf := range actorOfs {
		actorOf.sendSystemMessage(stopMessage)
	}
}

func (*actorContext) ResumeChildren(actorOfs ...*ActorRef) {
	for _, actorOf := range actorOfs {
		actorOf.sendSystemMessage(resumeMailboxMessage)
	}
}

//
// Miscellaneous
//

func (ctx *actorContext) GoString() string {
	return ctx.self.String()
}

func (ctx *actorContext) String() string {
	return ctx.self.String()
}
