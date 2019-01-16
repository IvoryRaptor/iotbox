package mock

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/ahandler"
	"time"
)

type Mock struct {
	ahandler.AHandlers
	channel chan common.ITask
	packet  common.Packet
}

func (m *Mock) Run() {
	m.channel <- m
}

func (m *Mock) Work(channel common.IModule) {
	fmt.Printf("[mock] Work\n")
	var packet common.Packet
	for i := 0; i < 10 && packet == nil; i++ {
		ch := channel.Send(m.packet)
		select {
		case res := <-ch:
			packet = res
		case <-time.After(time.Second * 3):
			fmt.Println("Timeout!")
		}
		if packet != nil {
			break
		}
	}
	if packet != nil {
		fmt.Printf("[mock] %d Complete\n", packet["value"])
		m.WorkHandlers(packet)
	}
}

func (m *Mock) Config(kernel common.IKernel, config map[interface{}]interface{}) error {
	m.channel = kernel.GetModule(config["module"].(string))
	m.packet = config["packet"].(map[interface{}]interface{})
	m.ConfigHandlers(kernel, config["handler"].([]interface{}))
	return nil
}
