package common

import (
	"github.com/robfig/cron"
)

type Packet map[interface{}]interface{}

type IKernel interface {
	GetModule(name string) chan ITask
	JoinTask(spec string, task cron.Job)
	CreateModule(config map[string]interface{}) (IModule, error)
	CreateTask(config map[interface{}]interface{}) (ITask, error)
}
