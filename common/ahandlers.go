package common

type AHandlers struct {
	ATask
	handlers []IHandlerTask
}

func (m *AHandlers) ConfigHandlers(kernel IKernel, configs []interface{}) error {
	m.handlers = make([]IHandlerTask, len(configs))
	for i, c := range configs {
		config := c.(map[interface{}]interface{})
		if item, err := kernel.CreateTask(config); err != nil {
			return err
		} else {
			m.handlers[i] = item.(IHandlerTask)
		}
	}
	return nil
}

func (m *AHandlers) WorkHandlers(packet Packet) {
	for _, handler := range m.handlers {
		handler.Clone().SetPacket(packet).Run()
	}
}
