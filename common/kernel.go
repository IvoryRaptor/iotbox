package common

//通讯数据包
type Packet map[string]interface{}

//内核接口
type IKernel interface {
	//通过名称，获取模块
	GetModule(name string) chan ITask //返回模块通道
	//创建模块
	CreateModule(config map[string]interface{}) (IModule, error) //模块，错误
	//创建任务
	CreateTask(config map[string]interface{}) (ITask, error) //任务，错误
}
