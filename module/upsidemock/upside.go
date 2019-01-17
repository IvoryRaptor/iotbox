package upsidemock

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/owner"
	"github.com/robfig/cron"
)

type Upside struct {
	common.AModule
	cron *cron.Cron
}

func (m *Upside) Config(kernel common.IKernel, config map[string]interface{}) error {
	m.cron = cron.New()
	task := &owner.Owner{}
	task.Config(
		kernel,
		map[interface{}]interface{}{
			"target": "downsidemock",
			"packet": map[interface{}]interface{}{
				"address": "bbb",
			},
		},
	)
	task.SetOwner(m)
	m.cron.AddJob("*/5 * * * * ?", task)
	m.cron.Start()
	return nil
}

func (m *Upside) Send(task common.ITask, packet common.Packet) chan common.Packet {
	fmt.Printf("Owner Complete Value = %d\n", packet["value"])
	m.Response <- nil
	return m.Response
}
