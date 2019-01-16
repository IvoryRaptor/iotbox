package common

type IModule interface {
	Config(ch chan ITask, config map[string]interface{}) error
	Send(packet Packet) chan Packet
}

type AModule struct {
	Response chan Packet
}

func (m *AModule) Start(ch chan ITask, this IModule) {
	m.Response = make(chan Packet)
	go func() {
		for {
			task := <-ch
			task.Work(this)
		}
	}()
}
