package common

type Packet map[string]interface{}

type IKernel interface {
	GetModule(name string) chan ITask
	CreateModule(config map[string]interface{}) (IModule, error)
	CreateTask(config map[string]interface{}) (ITask, error)
}
