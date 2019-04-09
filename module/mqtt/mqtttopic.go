package mqtt

import "github.com/IvoryRaptor/iotbox/akka"

type MQTTTopicActor struct {
	akka.Actor
}

func (actor *MQTTTopicActor) PreStart() error {
	return nil
}
func (actor *MQTTTopicActor) Config(config map[string]interface{}) error {
	return nil
}

func (actor *MQTTTopicActor) Receive(sender akka.IActor, message akka.Message) error {
	return nil
}
