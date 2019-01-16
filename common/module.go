package common

type IModule interface {
	Config(config map[string]interface{}) error
	Send(packet Packet) chan Packet
	Start(ch chan ITask, this IModule)
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
