package mocktask

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type MockTask struct {
	channel chan common.ITask
	packet  common.Packet
}

func (m *MockTask) Run() {
	m.channel <- m
}

func (m *MockTask) Work(channel common.IChannel) {
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
	}
}

func (m *MockTask) Config(kernel common.IKernel, config map[string]interface{}) error {
	m.channel = kernel.GetChannel(config["channel"].(string))
	m.packet = config["packet"].(map[interface{}]interface{})
	return nil
}
