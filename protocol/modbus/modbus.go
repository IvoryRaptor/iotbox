package modbus

import (
	"encoding/binary"
	"log"
)

// Protocol modbus 协议
type Protocol struct {
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
	itemType string
	// 功能码
	funcCode byte
}

type dataUnit struct {
	funcCode byte
	data     []byte
}

// Config 配置协议解析和编码的字段
func (mp *Protocol) Config(config map[interface{}]interface{}) (err error) {
	// 传输标识（net使用）
	log.Println(config)
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
	if _, ok := config["itemType"]; ok {
		mp.itemType = config["itemType"].(string)
	}
	log.Println("address", mp.registerAddress , "len", mp.registerLen)
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
