package main

import (
	"fmt"
	"log"
)

// RTUModbusProtocol modbus 协议
type RTUModbusProtocol struct {
	Protocol
}

// CreateRTUModbusProtocol 构造器
func CreateRTUModbusProtocol() *RTUModbusProtocol {
	return &RTUModbusProtocol{}
}

// GetName 获取协议名
func (mp *RTUModbusProtocol) GetName() string {
	return "RTUModbus"
}

// Encode 组包
func (mp *RTUModbusProtocol) Encode(config map[string]interface{}) (data []byte, err error) {
	dataUnit, err := mp.byteEncode(config)
	if err != nil {
		return nil, err
	}
	return mp.packager(dataUnit)
}

// Decode 解包
func (mp *RTUModbusProtocol) Decode(data []byte) (res map[string]interface{}, err error) {
	itemData := data[3 : len(data)-2]
	log.Printf("[%s]===> Decode data % X\n", mp.GetName(), itemData)
	return mp.byteDecode(itemData)
}

// Verify 包校验
func (mp *RTUModbusProtocol) Verify(data []byte) (err error) {
	length := len(data)
	if len(data) < 6 {
		return fmt.Errorf("[%s]==> len error[%d]", mp.GetName(), length)
	}
	var crc crc
	crc.reset().pushBytes(data[0 : length-2])
	checksum := uint16(data[length-1])<<8 | uint16(data[length-2])
	if checksum != crc.value() {
		return fmt.Errorf("[%s]==> crc error[%X] [%X]", mp.GetName(), checksum, crc.value())
	}
	// 设备地址
	if data[0] != mp.slaveID {
		err := fmt.Errorf("[%s]==> protocol id not match send[%x] recv[%x]",
			mp.GetName(), mp.slaveID, data[6])
		return err
	}
	// 功能码
	if data[1] != mp.funcCode {
		err := fmt.Errorf("[%s]==> protocol id not match send[%x] recv[%x]",
			mp.GetName(), mp.funcCode, data[7])
		return err
	}
	return nil
}

// 主机读
// 1byte 从地址
// 1byte 功能码
// 2byte 寄存器地址
// 1byte 寄存器个数
// 2byte crc16

// 主机写
// 1byte 从地址
// 1byte 功能码
// 2byte 寄存器地址
// 2byte 数据
// 2byte crc16

// packager 打包器
func (mp *RTUModbusProtocol) packager(pdu *dataUnit) (adu []byte, err error) {
	adu = make([]byte, 4+len(pdu.data))
	// 从地址
	adu[0] = mp.slaveID
	// 功能码
	adu[1] = mp.funcCode
	// 数据
	copy(adu[2:], pdu.data)
	// 加入crc
	var crc crc
	crc.reset().pushBytes(adu[0 : len(adu)-2])
	checksum := crc.value()
	adu[len(adu)-1] = byte(checksum >> 8)
	adu[len(adu)-2] = byte(checksum)
	return
}
