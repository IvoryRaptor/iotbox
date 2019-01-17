package mock

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type Mock struct {
	common.AHandlers
	packet common.Packet
}

func (m *Mock) MockWork(channel common.IModule) {
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

func (m *Mock) MockConfig(kernel common.IKernel, config map[interface{}]interface{}) error {
	m.packet = config["packet"].(map[interface{}]interface{})
	m.ConfigHandlers(kernel, config["handler"].([]interface{}))
	return nil
}

func Create() *Mock {
	result := &Mock{}
	result.SetCurrentWork(result.MockWork).SetOtherConfig(result.MockConfig)
	return result
}
