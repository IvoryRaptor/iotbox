package modbus

import (
	"encoding/binary"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
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
func (mp *NetModbusProtocol) Encode(config map[interface{}]interface{}) (data []byte, err error) {
	dataUnit := dataUnit{}
	dataUnit.funcCode = mp.funcCode
	switch mp.funcCode {
	case 0x01, 0x02, 0x03, 0x04:
		{
			dataUnit.data = dataBlock(mp.registerAddress, mp.registerLen)
		}
	default:
		return nil, fmt.Errorf("%s funcCode error[%x]", mp.GetName(), mp.funcCode)
	}
	return mp.packager(&dataUnit)
}

// Decode 解包
func (mp *NetModbusProtocol) Decode(data []byte) (item []common.ADataItem, err error) {
	return nil, nil
}

// Verify 包校验
func (mp *NetModbusProtocol) Verify(data []byte) (err error) {
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
