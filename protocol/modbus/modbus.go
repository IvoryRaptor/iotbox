package main

import (
	"encoding/binary"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/protocol"
	"github.com/fatih/structs"
	"log"
	"strings"
	"time"
)

// Protocol modbus 协议
type Protocol struct {
	protocol.AProtocol
	// 数据项目名称，全局不能重复 和ADataItem Name 进行bind
	name string
	// 传输标识（net使用）
	transactionID uint16
	// 设备地址
	slaveID byte
	// 寄存地址
	registerAddress uint16
	// 寄存器长度
	registerLen uint16
	// 是否是写操作
	isWrite bool
	// bool int string float
	valueType string
	// 功能码
	funcCode byte
	timeout  int
}

type dataUnit struct {
	funcCode byte
	data     []byte
}

// Create 构造器
func Create(config interface{}) (p protocol.IProtocol, err error) {
	mType, ok := config.(string)
	if !ok {
		return nil, fmt.Errorf("type not string")
	}
	switch strings.ToLower(mType) {
	case "net":
		return CreateNetModbusProtocol(), nil
	case "ascii":
		return nil, fmt.Errorf("ascii not support")
	case "rtu":
		return CreateRTUModbusProtocol(), nil
	}
	return nil, fmt.Errorf("type[%s] not support", mType)
}

// Config 配置协议解析和编码的字段
func (mp *Protocol) Config(config map[string]interface{}) (err error) {
	// 传输标识（net使用）
	log.Println(config)

	if _, ok := config["name"]; ok {
		mp.name = config["name"].(string)
	}

	if _, ok := config["transactionID"]; ok {
		mp.transactionID = uint16(config["transactionID"].(int))
	} else {
		mp.transactionID = 0x00
	}
	// 站地址
	if _, ok := config["slaveID"]; ok {
		mp.slaveID = byte(config["slaveID"].(int))
	}
	// 寄存器地址
	if _, ok := config["address"]; ok {
		mp.registerAddress = uint16(config["address"].(int))
	}
	// 寄存器长度
	if _, ok := config["len"]; ok {
		mp.registerLen = uint16(config["len"].(int))
	}
	// 是否是写
	if _, ok := config["isWrite"]; ok {
		mp.isWrite = config["isWrite"].(bool)
	}
	// 数据类型
	if _, ok := config["valueType"]; ok {
		mp.valueType = config["valueType"].(string)
	}
	// 超时时间
	if _, ok := config["timeout"]; ok {
		mp.timeout = config["timeout"].(int)
	}

	log.Println("address", mp.registerAddress, "len", mp.registerLen)
	if mp.isWrite {
		switch {
		case mp.registerAddress >= 1 && mp.registerAddress <= 9999:
			{
				// 线圈
				if mp.registerLen == 1 {
					mp.funcCode = 0x05
				} else {
					mp.funcCode = 0x0F
				}
				mp.registerAddress--
			}
		case mp.registerAddress >= 10001 && mp.registerAddress <= 19999:
			{
				// 离散
				if mp.registerLen == 1 {
					mp.funcCode = 0x06
				} else {
					mp.funcCode = 0x10
				}
				mp.registerAddress -= 10001
			}
		case mp.registerAddress >= 30001 && mp.registerAddress <= 39999:
			{
				// 输入寄存器 不能写
			}
		case mp.registerAddress >= 40001 && mp.registerAddress <= 49999:
			{
				// 保存寄存器
				if mp.registerLen == 1 {
					mp.funcCode = 0x06
				} else {
					mp.funcCode = 0x10
				}
				mp.registerAddress -= 40001
			}
		}
	} else {
		switch {
		case mp.registerAddress >= 1 && mp.registerAddress <= 9999:
			{
				// 线圈
				mp.funcCode = 0x01
				mp.registerAddress--
			}
		case mp.registerAddress >= 10001 && mp.registerAddress <= 19999:
			{
				// 离散
				mp.funcCode = 0x02
				mp.registerAddress -= 10001
			}
		case mp.registerAddress >= 30001 && mp.registerAddress <= 39999:
			{
				// 输入寄存器
				mp.funcCode = 0x04
				mp.registerAddress -= 30001
			}
		case mp.registerAddress >= 40001 && mp.registerAddress <= 49999:
			{
				// 保存寄存器
				mp.funcCode = 0x03
				mp.registerAddress -= 40001
			}
		}
	}
	return nil
}

// Decode 解包
func (mp *Protocol) byteDecode(data []byte) (res map[string]interface{}, err error) {
	item := common.ADataItem{Name: mp.name, ValueType: mp.valueType, SampleTime: time.Now()}

	switch mp.valueType {
	case "int", "float":
		// 可以对字节任意排序
		if len(data) < 4 {
			data = append(make([]byte, 4-len(data)), data...)
		}
	}
	if err := item.ByteToValue(data); err != nil {
		log.Printf("[modbus] ===> Decode ByteToValue %s\n", err)
	}
	log.Printf("[modbus] ===> %#v\n", item)
	// 对于int和float可以进行数据转换，是否有必要对转换公式进行抽象
	return common.Packet{
		"type":   "factors",
		"status": "ok",
		"value":  []map[string]interface{}{structs.Map(item)},
	}, nil
}

func (mp *Protocol) byteEncode(config map[string]interface{}) (res *dataUnit, err error) {
	dataUnit := dataUnit{}
	dataUnit.funcCode = mp.funcCode
	switch mp.funcCode {
	case 0x01, 0x02, 0x03, 0x04:
		{
			dataUnit.data = dataBlock(mp.registerAddress, mp.registerLen)
		}
	default:
		return nil, fmt.Errorf("[modbus] ==> funcCode error[%x]", mp.funcCode)
	}
	return &dataUnit, nil
}

// dataBlock creates a sequence of uint16 data.
func dataBlock(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	return data
}

// dataBlockSuffix creates a sequence of uint16 data and append the suffix plus its length.
func dataBlockSuffix(suffix []byte, value ...uint16) []byte {
	length := 2 * len(value)
	data := make([]byte, length+1+len(suffix))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	data[length] = uint8(len(suffix))
	copy(data[length+1:], suffix)
	return data
}
