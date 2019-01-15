package common

type IChannel interface {
	Config(ch chan ITask, config map[string]interface{}) error
	Send(packet Packet) chan Packet
}
