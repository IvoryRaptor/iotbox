package ahandler

import (
	"github.com/IvoryRaptor/iotbox/common"
)

type AHandlers struct {
	handlers []common.IHandlerTask
}

func (m *AHandlers) ConfigHandlers(kernel common.IKernel, configs []interface{}) error {
	m.handlers = make([]common.IHandlerTask, len(configs))
	for i, c := range configs {
		config := c.(map[interface{}]interface{})
		if item, err := kernel.CreateTask(config); err != nil {
			return err
		} else {
			m.handlers[i] = item.(common.IHandlerTask)
		}
	}
	return nil
}

func (m *AHandlers) WorkHandlers(packet common.Packet) {
	for _, handler := range m.handlers {
		handler.Clone().SetPacket(packet).Run()
	}
}
