package common

type Packet map[interface{}]interface{}

type IKernel interface {
	GetChannel(name string) chan ITask
	JoinTask(spec string, task ITask)
}
