package akka

type InitActor struct {
	Actor
}

func (actor *InitActor) Receive(sender IActor, message Message) error {
	return nil
}

func (actor *InitActor) Config(config map[string]interface{}) error {
	return nil
}
