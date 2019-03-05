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
		return nil, fmt.Errorf("[rootcloud]==> Encode not find type")
	}
	if _, ok := config["value"]; !ok {
		return nil, fmt.Errorf("[rootcloud]==> Encode not find value")
	}
	m := make(map[string]interface{})
	switch cType {
	case "factor":
		item := config["value"].(common.ADataItem)
		m[item.Name] = item.ConversionValue
		res, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("[rootcloud]==> Encode json[%s]", err)
		}
		return res, nil
	case "factors":
		for _, item := range config["value"].([]common.ADataItem) {
			m[item.Name] = item.ConversionValue
		}
		res, err := json.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("[rootcloud]==> Encode json[%s]", err)
		}
		return res, nil
	default:
		return nil, fmt.Errorf("[rootcloud]==> Encode type[%s] error", cType)
	}
}

// Decode 解码协议
func (p *Protocol) Decode(data []byte) (res map[string]interface{}, err error) {
	return nil, nil
}
