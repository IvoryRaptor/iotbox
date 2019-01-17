package common

type Packet map[interface{}]interface{}

type IKernel interface {
	GetModule(name string) chan ITask
	CreateModule(config map[string]interface{}) (IModule, error)
	CreateTask(config map[interface{}]interface{}) (ITask, error)
}
