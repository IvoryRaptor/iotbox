package common

// IProtocol 协议标准接口
type IProtocol interface {
	GetName() string
	Encode(config map[string]interface{}) (data []byte, err error)
	Config(config map[string]interface{}) ( err error)
	Decode(data []byte) (res map[string]interface{}, err error)
	Verify(data []byte) (err error)
}

// AProtocol 协议结构体
type AProtocol struct {
}
