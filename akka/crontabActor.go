package akka

//计划任务Actor

type CrontabActor struct {
	Actor
}

func (actor *CrontabActor) PreStart() error {
	return nil
}

func (actor *CrontabActor) Receive(sender IActor, message Message) error {
	return nil
}
