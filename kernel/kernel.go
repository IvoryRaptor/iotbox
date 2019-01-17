package kernel

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
)

type Kernel struct {
	channel map[string]chan common.ITask
}

func (k *Kernel) Start() {
	fmt.Printf("System Start\n")
}

func (k *Kernel) GetModule(name string) chan common.ITask {
	return k.channel[name]
}
