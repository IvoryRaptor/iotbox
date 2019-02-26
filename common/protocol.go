package common

// IProtocol 协议标准接口
type IProtocol interface {
	GetName() string
	Encode(config map[interface{}]interface{}) (data []byte, err error)
	Config(config map[interface{}]interface{}) ( err error)
	Decode(data []byte) (items []ADataItem, err error)
	Verify(data []byte) (err error)
}

// AProtocol 协议结构体
type AProtocol struct {
	// 数据项目名称，全局不能重复 和ADataItem Name 进行bind
	Name string
}
