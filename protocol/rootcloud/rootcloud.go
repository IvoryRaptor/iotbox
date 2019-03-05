package rootcloud

import (
	"encoding/json"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
)

// Protocol 根云上报协议
type Protocol struct {
	common.AProtocol
}

// GetName 获取协议名
func (p *Protocol) GetName() string {
	return "NetModbus"
}

// Verify 校验
func (p *Protocol) Verify(data []byte) (err error) {
	return nil
}

// Config 配置协议
func (p *Protocol) Config(config map[string]interface{}) (err error) {
	return nil
}

// Encode 编码
func (p *Protocol) Encode(config map[string]interface{}) (data []byte, err error) {
	var cType string
	if _, ok := config["type"]; ok {
		cType = config["type"].(string)
	} else {
		return nil, fmt.Errorf("[%s]==> Encode not find type", p.GetName())
	}
	if _, ok := config["value"]; !ok {
		return nil, fmt.Errorf("[%s]==> Encode not find value", p.GetName())
	}
	m := make(map[string]interface{})
	switch cType {
	case "factor":
		item := config["value"].(common.ADataItem)
		m[item.Name] = item.ConversionValue
		res, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("[%s]==> Encode json[%s]", err, p.GetName())
		}
		return res, nil
	case "factors":
		for _, item := range config["value"].([]common.ADataItem) {
			m[item.Name] = item.ConversionValue
		}
		res, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("[%s]==> Encode json[%s]", err, p.GetName())
		}
		return res, nil
	default:
		return nil, fmt.Errorf("[%s]==> Encode type[%s] error", cType, p.GetName())
	}
}

// Decode 解码协议
func (p *Protocol) Decode(data []byte) (res map[string]interface{}, err error) {
	return nil, nil
}

// Create 构造器
func Create() *Protocol {
	return &Protocol{}
}
