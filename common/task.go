package common

import "github.com/robfig/cron"

type ITask interface {
	cron.Job
	Config(kernel IKernel, config map[interface{}]interface{}) error
	Work(channel IModule)
}

type IHandlerTask interface {
	ITask
	Clone() IHandlerTask
	SetPacket(packet Packet) IHandlerTask
}

type IOwnerTask interface {
	ITask
	GetRequest() Packet
	SetOwner(owner IModule)
}
