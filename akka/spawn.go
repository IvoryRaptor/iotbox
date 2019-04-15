package akka

// ActorOf starts a new actor based on props and named with a unique id
// Deprecated: Use context.ActorOf instead.
func ActorOf(props *Props) *ActorRef {
	return EmptyRootContext.ActorOf(props)
}

// ActorOfPrefix starts a new actor based on props and named using a prefix followed by a unique id
// Deprecated: Use context.ActorOfPrefix instead.
func ActorOfPrefix(props *Props, prefix string) *ActorRef {
	return EmptyRootContext.ActorOfPrefix(props, prefix)
}

// ActorOfNamed starts a new actor based on props and named using the specified name
//
// If name exists, error will be ErrNameExists
// Deprecated: Use context.ActorOfNamed instead.
func ActorOfNamed(props *Props, name string) (*ActorRef, error) {
	context := EmptyRootContext
	return context.ActorOfNamed(props, name)
}
