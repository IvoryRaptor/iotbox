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

func (m *Upside) GetName() string {
	return "upside"
}

func (m *Upside) Config(kernel common.IKernel, config map[string]interface{}) error {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			task := owner.CreateOwner()
			task.Config(
				kernel,
				map[interface{}]interface{}{
					"target": "downsidemock",
					"request": map[interface{}]interface{}{
						"address": "ddd",
					},
				},
			)
			task.SetOwner(m).Run()
		}
	}()
	return nil
}

func (m *Upside) Send(task common.ITask, packet common.Packet) chan common.Packet {
	fmt.Printf("upside send value = %d\n", packet["value"])
	return m.Response
}

func CreateUpside() *Upside {
	return &Upside{}
}
