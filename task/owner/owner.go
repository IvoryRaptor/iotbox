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
	owner    common.IModule
	request  common.Packet
	response common.Packet
}

func (m *Owner) GetRequest() common.Packet {
	return m.response
}

func (m *Owner) WorkTarget(module common.IModule) (common.WorkState, error) {
	fmt.Printf("[owner] WorkTarget\n")
	for i := 0; i < 10 && m.response == nil; i++ {
		ch := module.Send(m, m.request)
		if m.response = module.Read(ch, time.Second*3); m.response != nil {
			break
		}
		fmt.Println("Timeout!")
	}
	m.SetCurrentWork(m.WorkOwner)
	m.owner.GetTaskQueue() <- m
	return common.Complete, nil
}

func (m *Owner) WorkOwner(channel common.IModule) (common.WorkState, error) {
	fmt.Printf("[owner] WorkOwner\n")
	channel.Send(m, m.response)
	return common.Complete, nil
}

func (m *Owner) SetOwner(owner common.IModule) *Owner {
	m.owner = owner
	return m
}

func (m *Owner) OwnerConfig(kernel common.IKernel, config map[interface{}]interface{}) error {
	m.request = config["request"].(map[interface{}]interface{})
	return nil
}

func Create() *Owner {
	result := &Owner{}
	result.SetCurrentWork(result.WorkTarget).SetOtherConfig(result.OtherConfig)
	return result
}
