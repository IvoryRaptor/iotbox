package mock

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/ahandler"
	"time"
)

type Mock struct {
	ahandler.AHandlers
	packet common.Packet
	target string
	kernel common.IKernel
}

func (m *Mock) Work(channel common.IModule) {
	fmt.Printf("[mock] Work\n")
	var packet common.Packet
	for i := 0; i < 10 && packet == nil; i++ {
		ch := channel.Send(m, m.packet)
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
func (m *Mock) Run() {
	m.kernel.GetModule(m.target) <- m
}
func (m *Mock) Config(kernel common.IKernel, config map[interface{}]interface{}) error {
	m.kernel = kernel
	m.target = config["target"].(string)
	m.packet = config["packet"].(map[interface{}]interface{})
	m.ConfigHandlers(kernel, config["handler"].([]interface{}))
	return nil
}
