package demo

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type Demo struct {
	common.AHandlers
	packet common.Packet
}

func (d *Demo) DemoWork(channel common.IModule) (common.WorkState, error) {
	fmt.Printf("[demo] Work\n")
	var packet common.Packet
	for i := 0; i < 10 && packet == nil; i++ {
		ch := channel.Send(d, d.packet)
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
		fmt.Printf("[demo] %d Complete\n", packet["value"])
		d.WorkHandlers(packet)
	}
	return common.Complete, nil
}

func (d *Demo) DemoConfig(kernel common.IKernel, config map[interface{}]interface{}) error {
	d.packet = config["packet"].(map[interface{}]interface{})
	d.ConfigHandlers(kernel, config["handler"].([]interface{}))
	return nil
}

func Create() *Demo {
	result := &Demo{}
	result.SetCurrentWork(result.DemoWork).SetOtherConfig(result.DemoConfig)
	return result
}
