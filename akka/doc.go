/*
Package actor declares the types used to represent actors in the Actor Model.

The actors model provide a high level abstraction for writing concurrent and distributed systems. This approach
simplifies the burden imposed on engineers, such as explicit locks and concurrent access to shared state, as actors
receive messages synchronously.

The following quote from Wikipedia distills the definition of an actor down to its essence

	In response to a message that it receives, an actor can: make local decisions, create more actors,
	send more messages, and determine how to respond to the next message received.


Creating Actors

Props provide the building blocks for declaring how actors should be created. The following example defines an actor
using a function literal to process messages:

	var props Props = actor.PropsFromFunc(func(c Context) {
		// process messages
	})

Alternatively, a type which conforms to the Actor interface, by defining a single Receive method, can be used.

	type MyActor struct {}

	func (a *MyActor) Receive(c Context) {
		// process messages
	}

	var props Props = actor.PropsFromProducer(func() Actor { return &MyActor{} })

ActorOf and ActorOfNamed use the given props to create a running instances of an actor. Once actorOfed, the actor is
ready to process incoming messages. To actorOf an actor with a unique name, use

	 actorOf := context.ActorOf(props)

The result of calling ActorOf is a unique ActorRef or process identifier.

Each time an actor is actorOfed, a new mailbox is created and associated with the ActorRef. Messages are sent to the mailbox
and then forwarded to the actor to process.


Processing Messages

An actor processes messages via its Receive handler. The signature of this function is:

	Receive(c actor.Context)

The actor system guarantees that this method is called synchronously, therefore there is no requirement to protect
shared state inside calls to this function.

Communicating With Actors

A ActorRef is the primary interface for sending messages to actors. The ActorRef.Tell method is used to send an asynchronous
message to the actor associated with the ActorRef:

	 actorOf.Tell("Hello World")

Depending on the requirements, communication between actors can take place synchronously or asynchronously. Regardless
of the circumstances, actors always communicate via a ActorRef.

When sending a message using ActorRef.Request or ActorRef.Ask, the actor which receives the message will respond
using the Context.Sender method, which returns the ActorRef of of the sender.

For synchronous communication, an actor will use a Future and wait for the result before continuing. To send a message
to an actor and wait for a response, use the Ask method, which returns a Future:

	f := actor.Ask( actorOf,"Hello", 50 * time.Millisecond)
	res, err := f.Result() // waits for  actorOf to reply */
package akka
