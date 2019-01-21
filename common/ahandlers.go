package common

type AHandlers struct {
	ATask
	handlers []ICloneTask
}

func (m *AHandlers) ConfigHandlers(kernel IKernel, configs []interface{}) error {
	m.handlers = make([]ICloneTask, len(configs))
	for i, c := range configs {
		config := c.(map[interface{}]interface{})
		if item, err := kernel.CreateTask(config); err != nil {
			return err
		} else {
			m.handlers[i] = item.(ICloneTask)
		}
	}
	return nil
}

func (m *AHandlers) WorkHandlers(packet Packet) {
	for _, handler := range m.handlers {
		if task, err := handler.Clone().SetPacket(packet); err == nil {
			task.Run()
		}
	}
}
