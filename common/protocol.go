package common

// IProtocol 协议标准接口
type IProtocol interface {
	GetName() string
	Encode(config map[interface{}]interface{}) (data []byte, err error)
	Config(config map[interface{}]interface{}) ( err error)
	Decode(data []byte) (item []ADataItem, err error)
	Verify(data []byte) (err error)
}

// // AProtocol 协议结构体
// type AProtocol struct {
// 	// 编码时存入，方便解码或者校验使用
// 	config map[string]interface{}
// }
