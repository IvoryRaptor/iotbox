package ahandler

import (
	"errors"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/sql"
)

type AHandlers struct {
	common.ATask
	handlers []common.IHandlerTask
}

func (m *AHandlers) ConfigHandlers(kernel common.IKernel, configs []interface{}) error {
	m.handlers = make([]common.IHandlerTask, len(configs))
	for i, c := range configs {
		config := c.(map[interface{}]interface{})
		handlerType := config["type"].(string)
		var item common.IHandlerTask
		switch handlerType {
		case "sql":
			item = &sql.Sql{}
		}
		if item == nil {
			return errors.New(fmt.Sprintf("Unknown Handler Type [%s]", handlerType))
		}
		item.Config(kernel, config)
		fmt.Printf("Add Handler Task %s\n", handlerType)
		m.handlers[i] = item
	}
	return nil
}

func (m *AHandlers) WorkHandlers(packet common.Packet) {
	for _, handler := range m.handlers {
		handler.Clone().SetPacket(packet).Run()
	}
}
