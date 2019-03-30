package main

import (
	"encoding/binary"
	"fmt"
	"log"
)

const (
	tcpProtocolIdentifier uint16 = 0x0000
	tcpHeaderSize                = 7
	tcpMaxLength                 = 260
)

// NetModbusProtocol modbus 协议
type NetModbusProtocol struct {
	Protocol
}

// CreateNetModbusProtocol 构造器
func CreateNetModbusProtocol() *NetModbusProtocol {
	return &NetModbusProtocol{}
}

// GetName 获取协议名
func (mp *NetModbusProtocol) GetName() string {
	return "NetModbus"
}

// Encode 组包
func (mp *NetModbusProtocol) Encode(config map[string]interface{}) (data []byte, err error) {
	dataUnit, err := mp.byteEncode(config)
	if err != nil {
		return nil, err
	}
	return mp.packager(dataUnit)
}

// Decode 解包
func (mp *NetModbusProtocol) Decode(data []byte) (res map[string]interface{}, err error) {
	// 5个0 5byte
	// 长度 1byte
	// 站地址 1byte
	// 功能吗 1byte
	// 剩余长度 1byte
	// 数据
	itemData := data[9:]
	log.Printf("[%s]===> Decode data % X\n", mp.GetName(), itemData)
	return mp.byteDecode(itemData)
}

// Verify 包校验
func (mp *NetModbusProtocol) Verify(data []byte) (err error) {
	if len(data) < 9 {
		return fmt.Errorf("[%s]==> len error[%d]", mp.GetName(), len(data))
	}
	// 传输ID
	responseTransactionID := binary.BigEndian.Uint16(data)
	if responseTransactionID != mp.transactionID {
		err := fmt.Errorf("[%s]==> transaction id not match send[%x] recv[%x]",
			mp.GetName(), mp.transactionID, responseTransactionID)
		return err
	}
	// 协议ID
	responseProtocolID := binary.BigEndian.Uint16(data[2:])
	if responseProtocolID != tcpProtocolIdentifier {
		err := fmt.Errorf("[%s]==> protocol id not match send[%x] recv[%x]",
			mp.GetName(), tcpProtocolIdentifier, responseProtocolID)
		return err
	}
	// 设备地址
	if data[6] != mp.slaveID {
		err := fmt.Errorf("[%s]==> protocol id not match send[%x] recv[%x]",
			mp.GetName(), mp.slaveID, data[6])
		return err
	}
	// 功能码
	if data[7] != mp.funcCode {
		err := fmt.Errorf("[%s]==> protocol id not match send[%x] recv[%x]",
			mp.GetName(), mp.funcCode, data[7])
		return err
	}
	return nil
}

// packager 打包器
func (mp *NetModbusProtocol) packager(pdu *dataUnit) (adu []byte, err error) {
	adu = make([]byte, tcpHeaderSize+1+len(pdu.data))
	// Transaction identifier
	binary.BigEndian.PutUint16(adu, uint16(mp.transactionID))
	// Protocol identifier
	binary.BigEndian.PutUint16(adu[2:], tcpProtocolIdentifier)
	// Length = sizeof(slaveID) + sizeof(funcCode) + Data
	length := uint16(1 + 1 + len(pdu.data))
	binary.BigEndian.PutUint16(adu[4:], length)
	// Unit identifier
	adu[6] = mp.slaveID
	// PDU
	adu[tcpHeaderSize] = pdu.funcCode
	copy(adu[tcpHeaderSize+1:], pdu.data)
	return
}