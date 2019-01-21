package demo

import (
	"errors"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type Demo struct {
	common.AHandlers
	retryCount int
	request    []common.Packet
	index      int
}

func (d *Demo) StartWork(module common.IModule) (common.WorkState, error) {
	var response common.Packet
	for i := 0; i < d.retryCount && response == nil; i++ {
		ch := module.Send(d, d.request[d.index])
		if response = module.Read(ch, time.Second*3); response == nil {
			fmt.Println("Timeout!")
		}
	}
	if response == nil {
		d.index = (d.index + 1) % len(d.request) //Jump next
		return common.Failed, errors.New("Timeout ")
	}
	d.WorkHandlers(response)
	fmt.Printf("[demo] %d Complete\n", response["value"])
	//避免占用时间过长
	if d.index++; d.index < len(d.request) {
		return common.Running, nil
	}
	d.index = 0
	return common.Complete, nil
}

func (d *Demo) DemoConfig(kernel common.IKernel, config map[interface{}]interface{}) error {
	d.index = 0
	d.SetCurrentWork(d.StartWork)
	p := config["packet"].([]interface{})
	d.request = make([]common.Packet, len(p))
	for index, item := range p {
		d.request[index] = item.(map[interface{}]interface{})
	}
	d.retryCount = config["retry"].(int)
	d.ConfigHandlers(kernel, config["handler"].([]interface{}))
	return nil
}

func CreateDemo() *Demo {
	result := &Demo{}
	result.SetOtherConfig(result.DemoConfig)
	return result
}
