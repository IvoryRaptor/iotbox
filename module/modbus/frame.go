package modbus

import (
	"encoding/binary"
	"fmt"
)

// IModbusFrame modbus帧接口
type IModbusFrame interface {
	// Bit access

	// ReadCoils reads from 1 to 2000 contiguous status of coils in a
	// remote device and returns coil status.
	ReadCoils(address, quantity uint16) (results []byte, err error)
	// ReadDiscreteInputs reads from 1 to 2000 contiguous status of
	// discrete inputs in a remote device and returns input status.
	ReadDiscreteInputs(address, quantity uint16) (results []byte, err error)
	// WriteSingleCoil write a single output to either ON or OFF in a
	// remote device and returns output value.
	WriteSingleCoil(address, value uint16) (results []byte, err error)
	// WriteMultipleCoils forces each coil in a sequence of coils to either
	// ON or OFF in a remote device and returns quantity of outputs.
	WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error)

	// 16-bit access

	// ReadInputRegisters reads from 1 to 125 contiguous input registers in
	// a remote device and returns input registers.
	ReadInputRegisters(address, quantity uint16) (results []byte, err error)
	// ReadHoldingRegisters reads the contents of a contiguous block of
	// holding registers in a remote device and returns register value.
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
	// WriteSingleRegister writes a single holding register in a remote
	// device and returns register value.
	WriteSingleRegister(address, value uint16) (results []byte, err error)
	// WriteMultipleRegisters writes a block of contiguous registers
	// (1 to 123 registers) in a remote device and returns quantity of
	// registers.
	WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error)
	// ReadWriteMultipleRegisters performs a combination of one read
	// operation and one write operation. It returns read registers value.
	ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error)
	// MaskWriteRegister modify the contents of a specified holding
	// register using a combination of an AND mask, an OR mask, and the
	// register's current contents. The function returns
	// AND-mask and OR-mask.
	MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error)
	//ReadFIFOQueue reads the contents of a First-In-First-Out (FIFO) queue
	// of register in a remote device and returns FIFO value register.
	ReadFIFOQueue(address uint16) (results []byte, err error)
}
// AModbusFrame modbus 实现
type AModbusFrame struct {
	packager    Packager
}

// NewModbusFrame 创建组帧器
func NewModbusFrame(pack Packager) *AModbusFrame {
	return &AModbusFrame{packager: pack}
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
// ReadCoils 读取coils
//  Function code         : 1 byte (0x01)
//  Starting address      : 2 bytes
//  Quantity of coils     : 2 bytes
func (mb *AModbusFrame) ReadCoils(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 2000)
		return nil, err
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadCoils,
		Data: dataBlock(address, quantity),
	}
	res, err:= mb.packager.Encode(&request)
	return res, err
}
// ReadDiscreteInputs 读取输入寄存器
//  Function code         : 1 byte (0x02)
//  Starting address      : 2 bytes
//  Quantity of inputs    : 2 bytes
func (mb *AModbusFrame) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 2000 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 2000)
		return nil,err
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadDiscreteInputs,
		Data:         dataBlock(address, quantity),
	}
	return mb.packager.Encode(&request)
}

// ReadHoldingRegisters  读取保存寄存器
//  Function code         : 1 byte (0x03)
//  Starting address      : 2 bytes
//  Quantity of registers : 2 bytes
func (mb *AModbusFrame) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return nil,err
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadHoldingRegisters,
		Data:         dataBlock(address, quantity),
	}
	return mb.packager.Encode(&request)
}

// ReadInputRegisters 读取输入保存寄存器
//  Function code         : 1 byte (0x04)
//  Starting address      : 2 bytes
//  Quantity of registers : 2 bytes
func (mb *AModbusFrame) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return nil,err
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadInputRegisters,
		Data:         dataBlock(address, quantity),
	}
	return mb.packager.Encode(&request)
}

// WriteSingleCoil 写单线圈
//  Function code         : 1 byte (0x05)
//  Output address        : 2 bytes
//  Output value          : 2 bytes
func (mb *AModbusFrame) WriteSingleCoil(address, value uint16) (results []byte, err error) {
	// The requested ON/OFF state can only be 0xFF00 and 0x0000
	if value != 0xFF00 && value != 0x0000 {
		err = fmt.Errorf("modbus: state '%v' must be either 0xFF00 (ON) or 0x0000 (OFF)", value)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteSingleCoil,
		Data:         dataBlock(address, value),
	}
	return mb.packager.Encode(&request)
}

// WriteSingleRegister 写多线圈
//  Function code         : 1 byte (0x06)
//  Register address      : 2 bytes
//  Register value        : 2 bytes
func (mb *AModbusFrame) WriteSingleRegister(address, value uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteSingleRegister,
		Data:         dataBlock(address, value),
	}
	return mb.packager.Encode(&request)
}

// WriteMultipleCoils 写多个线圈
//  Function code         : 1 byte (0x0F)
//  Starting address      : 2 bytes
//  Quantity of outputs   : 2 bytes
//  Byte count            : 1 byte
//  Outputs value         : N* bytes
func (mb *AModbusFrame) WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error) {
	if quantity < 1 || quantity > 1968 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 1968)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteMultipleCoils,
		Data:         dataBlockSuffix(value, address, quantity),
	}
	return mb.packager.Encode(&request)
}

// WriteMultipleRegisters 写多个寄存器
//  Function code         : 1 byte (0x10)
//  Starting address      : 2 bytes
//  Quantity of outputs   : 2 bytes
//  Byte count            : 1 byte
//  Registers value       : N* bytes
func (mb *AModbusFrame) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	if quantity < 1 || quantity > 123 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 123)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteMultipleRegisters,
		Data:         dataBlockSuffix(value, address, quantity),
	}
	return mb.packager.Encode(&request)
}


// MaskWriteRegister
//  Function code         : 1 byte (0x16)
//  Reference address     : 2 bytes
//  AND-mask              : 2 bytes
//  OR-mask               : 2 bytes

func (mb *AModbusFrame) MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeMaskWriteRegister,
		Data:         dataBlock(address, andMask, orMask),
	}
	return mb.packager.Encode(&request)
}

// ReadWriteMultipleRegisters
//  Function code         : 1 byte (0x17)
//  Read starting address : 2 bytes
//  Quantity to read      : 2 bytes
//  Write starting address: 2 bytes
//  Quantity to write     : 2 bytes
//  Write byte count      : 1 byte
//  Write registers value : N* bytes

func (mb *AModbusFrame) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	if readQuantity < 1 || readQuantity > 125 {
		err = fmt.Errorf("modbus: quantity to read '%v' must be between '%v' and '%v',", readQuantity, 1, 125)
		return
	}
	if writeQuantity < 1 || writeQuantity > 121 {
		err = fmt.Errorf("modbus: quantity to write '%v' must be between '%v' and '%v',", writeQuantity, 1, 121)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadWriteMultipleRegisters,
		Data:         dataBlockSuffix(value, readAddress, readQuantity, writeAddress, writeQuantity),
	}
	return mb.packager.Encode(&request)
}

// ReadFIFOQueue
//  Function code         : 1 byte (0x18)
//  FIFO pointer address  : 2 bytes
// Response:
//  Function code         : 1 byte (0x18)
//  Byte count            : 2 bytes
//  FIFO count            : 2 bytes
//  FIFO count            : 2 bytes (<=31)
//  FIFO value register   : Nx2 bytes
func (mb *AModbusFrame) ReadFIFOQueue(address uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadFIFOQueue,
		Data:         dataBlock(address),
	}
	return mb.packager.Encode(&request)
}