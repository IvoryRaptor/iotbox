package owner

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type State struct {
}

type Owner struct {
	target string
	kernel common.IKernel
	packet common.Packet
	work   func(channel common.IModule)
	owner  common.IModule
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
	m.work = m.WorkOwner
	m.owner.GetTaskQueue() <- m
}

func (m *Owner) WorkOwner(channel common.IModule) {
	fmt.Printf("[owner] WorkOwner\n")
	<-channel.Send(m, m.result)
}

func (m *Owner) Run() {
	m.kernel.GetModule(m.target) <- m
}

func (m *Owner) Work(channel common.IModule) {
	m.work(channel)
}

func (m *Owner) SetOwner(owner common.IModule) {
	m.owner = owner
}

func (m *Owner) Config(kernel common.IKernel, config map[interface{}]interface{}) error {
	m.kernel = kernel
	m.target = config["target"].(string)
	m.packet = config["packet"].(map[interface{}]interface{})
	m.work = m.WorkTarget
	return nil
}
