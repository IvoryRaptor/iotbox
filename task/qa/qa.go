package qa

// 一问一答模式

import (
	"errors"
	"log"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

func init() {
	log.Println("Question and answer task init")
}

// QA Question and answer
type QA struct {
	common.AHandlers
	retryCount int
	request    []common.Packet
	index      int
}

// StartWork QA work
func (d *QA) StartWork(module common.IModule) (common.WorkState, error) {
	var response common.Packet
	for i := 0; i < d.retryCount && response == nil; i++ {
		log.Println("==============>0")
		ch := module.Send(d, d.request[d.index])
		log.Println("==============>1")
		if response = module.Read(ch, time.Second*5); response == nil {
			log.Println("[QA] Timeout!")
			log.Println("==============>2")
		}
		log.Println("==============>3")
	}
	if response == nil {
		d.index = (d.index + 1) % len(d.request) //Jump next
		return common.Failed, errors.New("Timeout ")
	}
	d.WorkHandlers(response)
	log.Printf("[QA] %d Complete\n", response["value"])
	//避免占用时间过长
	if d.index++; d.index < len(d.request) {
		return common.Running, nil
	}
	d.index = 0
	return common.Complete, nil
}
// QAConfig QA config
func (d *QA) QAConfig(kernel common.IKernel, config map[string]interface{}) error {
	d.index = 0
	d.SetCurrentWork(d.StartWork)
	p := config["packet"].([]interface{})
	d.request = make([]common.Packet, len(p))
	for index, item := range p {
		d.request[index] = item.(map[string]interface{})
	}
	d.retryCount = config["retry"].(int)
	d.ConfigHandlers(kernel, config["handler"].([]interface{}))
	return nil
}
// Create 创建QA任务
func Create() *QA {
	result := &QA{}
	result.SetOtherConfig(result.QAConfig)
	return result
}
