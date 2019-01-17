package upsidemock

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/owner"
	"time"
)

type Upside struct {
	common.AModule
}

func (m *Upside) Config(kernel common.IKernel, config map[string]interface{}) error {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			task := owner.Create()
			task.Config(
				kernel,
				map[interface{}]interface{}{
					"target": "downsidemock",
					"packet": map[interface{}]interface{}{
						"address": "bbb",
					},
				},
			)
			task.SetOwner(m).Run()
		}
	}()
	return nil
}

func (m *Upside) Send(task common.ITask, packet common.Packet) chan common.Packet {
	fmt.Printf("Owner Complete Value = %d\n", packet["value"])
	return m.Response
}
