package common

// 有执行列表的任务
type AHandlers struct {
	ATask
	handlers []ICloneTask //执行成功后自动执行的子任务列表，子任务必须继承自ICloneTask，避免一些资源反复初始化
}

//配置执行列表
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

//执行列表内容
func (m *AHandlers) WorkHandlers(packet Packet) {
	for _, handler := range m.handlers {
		if task, err := handler.Clone().SetPacket(packet); err == nil {
			task.Run()
		}
	}
}
