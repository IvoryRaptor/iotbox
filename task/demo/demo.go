package demo

import (
	"errors"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type Demo struct {
	common.AHandlers
	packet common.Packet
}

func (d *Demo) DemoWork(module common.IModule) (common.WorkState, error) {
	fmt.Printf("[demo] Work\n")
	var packet common.Packet
	for i := 0; i < 10 && packet == nil; i++ {
		ch := module.Send(d, d.packet)
		if packet = module.Read(ch, time.Second*3); packet != nil {
			break
		} else {
			fmt.Println("Timeout!")
		}
	}
	if packet == nil {
		return common.Failed, errors.New("")
	}
	fmt.Printf("[demo] %d Complete\n", packet["value"])
	d.WorkHandlers(packet)
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
