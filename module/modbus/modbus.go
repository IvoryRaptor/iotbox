package modbus
import (
	"github.com/IvoryRaptor/iotbox/common"
	"time"
	"log"
	"github.com/tarm/serial"
	"io"
	"net"
)

// Modbus 模块
type Modbus struct {
	common.AModule
	config map[string]interface{}
	wait        time.Duration
	frame string
	port string
	portConfig string
	address int
}
// GetName 获取模块名
func (m *Modbus) GetName() string {
	return "modbus"
}
// Config 配置模块
func (m *Modbus) Config(_ common.IKernel, config map[string]interface{}) error {
	log.Print("Config\n")
	m.config = config
	m.wait = time.Duration(config["wait"].(int)) * time.Second
	m.frame = config["frame"].(string)
	m.port = config["port"].(string)
	m.portConfig = config["portConfig"].(string)
	m.address = config["address"].(int)
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
	return buf,nil
}

// Send 发送数据
func (m *Modbus) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	config := &serial.Config{Name: m.port, Baud: 115200, ReadTimeout: m.wait}
	log.Printf("[%s] send\n", m.GetName())
	port, err := serial.OpenPort(config)
	if err != nil {
			log.Fatal(err)
			return m.Response
	}
	conn, err := net.Dial("tcp", "192.168.1.234:502")
	if err != nil {
		log.Fatal(err)
		return m.Response
	}
	go func() {
		defer conn.Close()
		conn.SetReadDeadline(time.Now().Add(time.Second*3))
		frame := NewModbusFrame(NewTCPPackager(0x00,0x01))
		buf,_ := frame.ReadInputRegisters(0,1)
		log.Println("send tcp frame", buf)
		conn.Write(buf)
		buf, err := read(conn)
		if err != nil {
			log.Fatal("read error  ", err)
			return
		}
		log.Println("recv tcp frame", buf)
		// m.Response <- common.Packet{
		// 	"value": buf,
		// }
	}()
	go func() {
		defer port.Close()
		port.Write([]byte(`{"type":"ping"}`))
		buf, err := read(port)
		if err != nil {
			log.Fatal("read error  ", err)
			return
		}
		n := len(buf)
		log.Printf("[%s] recv %d %q", m.GetName(), n, buf[:n])
		m.Response <- common.Packet{
			"value": buf,
		}
	}()
	return m.Response
}
// Create 创建modbus对象
func Create() *Modbus {
	log.Println("CreateModbus")
	return &Modbus{}
}