package modbus

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/protocol"
	"github.com/tarm/serial"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
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
	// 通道空闲时间,用于接收数据后空闲多久认为帧结束 单位ms
	idleTime int
	// 发送延时,单位ms
	sendDelay int
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

	if _, ok := config["idleTime"]; ok {
		m.idleTime = config["idleTime"].(int)
	}
	if _, ok := config["sendDelay"]; ok {
		m.sendDelay = config["sendDelay"].(int)
	}
	log.Printf("[%s]==> Config %#v\n", m.GetName(), m)
	return nil
}

func read(reader io.Reader, idleTime int, timeout int) ([]byte, error) {
	buf := []byte{}
	ch := make(chan []byte, 10)
	timer := time.NewTimer(time.Duration(timeout) * time.Millisecond)
	go func() {
		for {
			tmp := make([]byte, 128)
			len, err := reader.Read(tmp)
			if len == 0 || err != nil {
				return
			}
			if len > 0 {
				ch <- tmp[:len]
			}
		}
	}()
	for {
		select {
		case <-timer.C:
			{
				timer.Stop()
				goto breakout
			}
		case tmp := <-ch:
			{
				buf = append(buf, tmp...)
				timer.Reset(time.Duration(idleTime) * time.Millisecond)
			}
		}
	}
breakout:
	return buf, nil
}

// createConnect 创建链接
func (m *Modbus) createConnect() (io.ReadWriteCloser, error) {
	var res io.ReadWriteCloser
	switch strings.ToLower(m.portType) {
	case "serial":
		config := &serial.Config{Name: m.port, Baud: 9600}
		serialConfig := strings.Split(m.portConfig, ",")
		if len(serialConfig) == 4 {
			if baud, err := strconv.Atoi(serialConfig[0]); err == nil {
				config.Baud = baud
			}
			if size, err := strconv.Atoi(serialConfig[1]); err == nil {
				config.Size = byte(size)
			}
			switch strings.ToLower(serialConfig[2]) {
			case "n":
				config.Parity = serial.ParityNone
			case "o":
				config.Parity = serial.ParityOdd
			case "e":
				config.Parity = serial.ParityEven
			case "m":
				config.Parity = serial.ParityMark
			case "s":
				config.Parity = serial.ParitySpace
			}
			switch serialConfig[3] {
			case "1":
				config.StopBits = serial.Stop1
			case "2":
				config.StopBits = serial.Stop2
			case "1.5":
				config.StopBits = serial.Stop1Half
			}
		}
		log.Printf("[%s][%s]===>config %#v", m.GetName(), m.port, config)
		port, err := serial.OpenPort(config)
		if err != nil {
			return nil, fmt.Errorf("connect[%s] error[%s]", m.port, err)
		}
		res = port
	case "net":
		conn, err := net.Dial(strings.ToLower(m.portConfig), m.port)
		if err != nil {
			return nil, fmt.Errorf("net[%s] error[%s]", m.port, err)
		}
		res = conn
	default:
		log.Printf("portType error[%s]\n", m.portType)
		return nil, fmt.Errorf("portType error[%s]", m.portType)
	}
	return res, nil
}

// Tell 发送数据
func (m *Modbus) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	log.Printf("[%s][%s]==> Tell\n", m.GetName(), m.port)
	go func() {
		var conn io.ReadWriteCloser
		var p protocol.IProtocol
		var err error
		var sendBuf, recvBuf []byte
		var value map[string]interface{}
		timeout := 2000
		// 超时时间
		if _, ok := packet["timeout"]; ok {
			timeout = packet["timeout"].(int)
		}
		conn, err = m.createConnect()
		if err != nil {
			goto breakout
		}
		defer conn.Close()
		p, err = protocol.CreateProtocol("modbus", m.protocolType)
		if err != nil {
			goto breakout
		}
		p.Config(packet)
		sendBuf, err = p.Encode(packet)
		if err != nil {
			goto breakout
		}
		log.Printf("[%s][%s]==> Tell frame [% X]\n", m.GetName(), m.port, sendBuf)
		// sleep sendDelay time
		time.Sleep(time.Millisecond * time.Duration(m.sendDelay))
		conn.Write(sendBuf)
		recvBuf, err = read(conn, m.idleTime, timeout)
		if err != nil {
			goto breakout
		}
		log.Printf("[%s][%s]==> recv frame [% X]\n", m.GetName(), m.port, recvBuf)
		if err := p.Verify(recvBuf); err != nil {
			goto breakout
		}
		value, err = p.Decode(recvBuf)
		if err != nil {
			goto breakout
		}
		log.Printf("[%s]===> Decode value [%#v]\n", p.GetName(), value)
		m.Response <- value
		return
	breakout:
		log.Printf("[%s][%s]====> %s", m.GetName(), m.port, err)
		m.Response <- nil
		return
	}()
	return m.Response
}

// Create 创建modbus对象
func Create() *Modbus {
	log.Println("CreateModbus")
	return &Modbus{}
}
