package modbus

import "github.com/IvoryRaptor/iotbox/akka"

type ModbusActor struct {
	akka.Actor
}

func (actor *ModbusActor) Receive(sender akka.IActor, message akka.Message) error {
	return nil
}

func (actor *ModbusActor) PreStart() error {
	return nil
}

func (actor *ModbusActor) Config(config map[string]interface{}) error {
	return nil
}
