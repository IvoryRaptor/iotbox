package protocol
import (
	"plugin"
)
// 协议库路径
const protocolPath = "./lib/protocol/"
// 协议库后缀
const protocolLibSuffix = ".so"

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

// CreateProtocol 创建协议处理器
func CreateProtocol(name string, config interface{}) (protocol IProtocol, err error) {
	var symbol plugin.Symbol
	var p *plugin.Plugin
	protocol = nil
	err = nil
	p, err = plugin.Open(protocolPath + name + protocolLibSuffix)
	if err != nil {
		return protocol, err
	}
	symbol, err = p.Lookup("Create")
	if err != nil {
		return protocol, err
	}
	protocol, err = symbol.(func(interface{})(IProtocol, error))(config)
	return protocol, err
}