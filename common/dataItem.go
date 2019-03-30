package common

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

const (
	// Sampling 采样中
	Sampling = iota
	// SampleSucess 采样成功
	SampleSucess
	// SampleFail 采样失败
	SampleFail
	// SampleHold 采样保持
	SampleHold
)

// ADataItem 数据项格式
type ADataItem struct {
	Name string
	// 值类型 bool int string float
	ValueType string
	// 原始值
	RawValue interface{}
	// 转换后值,只有int float 支持转换
	ConversionValue interface{}
	// 采样时间
	SampleTime time.Time
	// 采样状态
	SampleStatus int
}

// ToValue []byte转换标准value
func (item *ADataItem) ByteToValue(data []byte) error {
	switch item.ValueType {
	case "bool":
		if data[0] == 0 {
			item.RawValue = false
		} else {
			item.RawValue = true
		}
	case "int":
		item.RawValue = binary.BigEndian.Uint32(data)
	case "string":
		item.RawValue = string(data)
	case "float":
		item.RawValue = math.Float32frombits(binary.BigEndian.Uint32(data))
	default:
		return fmt.Errorf("ByteToValue unknow type")
	}
	item.ConversionValue = item.RawValue
	return nil
}

// StringToValue 字符串转标准value
func (item *ADataItem) StringToValue(data string) error {
	return nil
}
