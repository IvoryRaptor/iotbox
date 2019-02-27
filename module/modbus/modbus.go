package modbus

import (
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/protocol/modbus"
	"github.com/tarm/serial"
	"io"
	"log"
	"net"
	"time"
	"strings"
)

// Modbus 模块
type Modbus struct {
	common.AModule
	// 协议类型和模块有关 对于modbus取值有（rtu、net、ascii）
	protocolType string
	// 模块的通道类型 对于modbus取值有（serial、net）
	portType string
	// 模块的通道地址 例如（192.168.1.234:502、/dev/tty.usbserial）
	port string
	// 模块的通道配置 对应serial（115200,8,n,1）对应net（udp、tcp）
	portConfig string
	// 模块使用通道前延时，单位ms
	idle int
}

// GetName 获取模块名
func (m *Modbus) GetName() string {
	return "modbus"
}

// Config 配置模块
func (m *Modbus) Config(_ common.IKernel, config map[string]interface{}) error {
	if _, ok := config["protocolType"]; ok {
		m.protocolType = config["protocolType"].(string)
	}

	if _, ok := config["portType"]; ok {
		m.portType = config["portType"].(string)
	}

	if _, ok := config["port"]; ok {
		m.port = config["port"].(string)
	}

	if _, ok := config["portConfig"]; ok {
		m.portConfig = config["portConfig"].(string)
	}

	if _, ok := config["idle"]; ok {
		m.idle = config["idle"].(int)
	}
	log.Printf("[%s]==> Config %#v\n", m.GetName(), m)
	return nil
}

func read(reader io.Reader) ([]byte, error) {
	buf := []byte{}
	for true {
		tmp := make([]byte, 128)
		time.Sleep(time.Millisecond * 100)
		len, err := reader.Read(tmp)
		buf = append(buf, tmp[:len]...)
		if err != nil {
			break
		}
	}
	return buf, nil
}

// Send 发送数据
func (m *Modbus) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	log.Printf("[%s]==> Send\n", m.GetName())
	var rw io.ReadWriter
	var iProtocol common.IProtocol
	switch strings.ToLower(m.portType) {
	case "serial":
		config := &serial.Config{Name: m.port, Baud: 115200, ReadTimeout: time.Second}
		port, err := serial.OpenPort(config)
		if err != nil {
			log.Fatalf("[%s]==> connect[%s] error[%s]\n",m.GetName(),m.port, err)
			return m.Response
		}
		defer port.Close()
		rw = port
	case "net":
		conn, err := net.Dial(strings.ToLower(m.portConfig), m.port)
		if err != nil {
			log.Fatalf("[%s]==> connect[%s] error[%s]\n",m.GetName(),m.port, err)
			return m.Response
		}
		defer conn.Close()
		conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		rw = conn
	default:
		log.Fatalf("%s]==> Send portType error[%s]\n", m.GetName(), m.portType)
		return m.Response
	}
	switch strings.ToLower(m.protocolType){
		case "rtu":
		case "net":
			iProtocol = modbus.CreateNetModbusProtocol()
		case "ascii":
		default:
			log.Fatalf("[%s]==> Send protocolType error[%s]\n", m.GetName(), m.protocolType)
			return m.Response
	}
	iProtocol.Config(packet)
	sendBuf, err := iProtocol.Encode(packet)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[%s]==> Send frame [% X]\n", m.GetName(), sendBuf)
		rw.Write(sendBuf)
		readBuf, err := read(rw)
		if err != nil {
			log.Fatalf("[%s]==> read error[%s]\n", m.GetName(), err)
			return m.Response
		}
		log.Printf("[%s]==> recv frame [% X]\n", m.GetName(), readBuf)
		if err := iProtocol.Verify(readBuf) ; err != nil {
			log.Fatalf("[%s] ===> verify error[%s]\n", iProtocol.GetName(), err)
			return m.Response
		}
		value, errDecode := iProtocol.Decode(readBuf)
		if errDecode != nil {
			log.Printf("[%s] ===> Decode error[%s]\n", iProtocol.GetName(), errDecode)
			return m.Response
		}
		log.Printf("[%s] ===> Decode value [%#v]\n", iProtocol.GetName(), value)
		m.Response <- common.Packet{
			"value": value,
		}
	return m.Response
}

// Create 创建modbus对象
func Create() *Modbus {
	log.Println("CreateModbus")
	return &Modbus{}
}
