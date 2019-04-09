package mqtt

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

// AMQTTActor 上报模块
type MQTTMasterActor struct {
	akka.Actor
}

func (actor *MQTTMasterActor) PreStart() error {
	return nil
}
func (actor *MQTTMasterActor) Config(config map[string]interface{}) error {
	return nil
}

func (actor *MQTTMasterActor) Receive(sender akka.IActor, message akka.Message) error {
	return nil
}
