package common

import "github.com/robfig/cron"

type ITask interface {
	cron.Job
	Config(kernel IKernel, config map[interface{}]interface{}) error
	Work(channel IModule)
}

type ATask struct {
	target chan ITask
	this   ITask
}

func (s *ATask) InitTarget(kernel IKernel, config map[interface{}]interface{}, this ITask) {
	s.target = kernel.GetModule(config["target"].(string))
	s.this = this
}

func (s *ATask) Run() {
	s.target <- s.this
}

type IHandlerTask interface {
	ITask
	Clone() IHandlerTask
	SetPacket(packet Packet) IHandlerTask
}
