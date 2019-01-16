package kernel

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/robfig/cron"
)

type Kernel struct {
	channel map[string]chan common.ITask
	cron    *cron.Cron
}

func (k *Kernel) Start() {
	fmt.Printf("System Start\n")
	k.cron.Start()
}

func (k *Kernel) GetModule(name string) chan common.ITask {
	return k.channel[name]
}

func (k *Kernel) JoinTask(spec string, task cron.Job) {
	k.cron.AddJob(spec, task)
}
