package owner

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type State struct {
}

type Owner struct {
	common.ATask
	owner  common.IModule
	packet common.Packet
	result common.Packet
}

func (m *Owner) GetRequest() common.Packet {
	return m.result
}

func (m *Owner) WorkTarget(channel common.IModule) {
	fmt.Printf("[owner] WorkTarget\n")
	for i := 0; i < 10 && m.result == nil; i++ {
		ch := channel.Send(m, m.packet)
		select {
		case res := <-ch:
			m.result = res
		case <-time.After(time.Second * 3):
			fmt.Println("Timeout!")
		}
		if m.result != nil {
			break
		}
	}
	m.SetCurrentWork(m.WorkOwner)
	m.owner.GetTaskQueue() <- m
}

func (m *Owner) WorkOwner(channel common.IModule) {
	fmt.Printf("[owner] WorkOwner\n")
	channel.Send(m, m.result)
}

func (m *Owner) SetOwner(owner common.IModule) *Owner {
	m.owner = owner
	return m
}

func (m *Owner) OwnerConfig(kernel common.IKernel, config map[interface{}]interface{}) error {
	m.packet = config["packet"].(map[interface{}]interface{})
	return nil
}

func Create() *Owner {
	result := &Owner{}
	result.SetCurrentWork(result.WorkTarget).SetOtherConfig(result.OtherConfig)
	return result
}
