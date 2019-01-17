package common

type IModule interface {
	Config(kernel IKernel, config map[string]interface{}) error
	Send(handle ITask, packet Packet) chan Packet
	Start(this IModule)
	GetTaskQueue() chan ITask
}

type AModule struct {
	Response  chan Packet
	taskQueue chan ITask
}

func (m *AModule) GetTaskQueue() chan ITask {
	return m.taskQueue
}

func (m *AModule) Start(this IModule) {
	m.Response = make(chan Packet)
	m.taskQueue = make(chan ITask, 10)
	go func() {
		for {
			task := <-m.GetTaskQueue()
			task.Work(this)
		}
	}()
}
