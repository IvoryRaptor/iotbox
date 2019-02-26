package common

import (
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

// ToValue 转换值
func (item *ADataItem) ToValue([]byte) error {
	return nil
}
