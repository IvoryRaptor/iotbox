package common

import "github.com/robfig/cron"

type Packet map[interface{}]interface{}

type IKernel interface {
	GetModule(name string) chan ITask
	JoinTask(spec string, task cron.Job)
}
