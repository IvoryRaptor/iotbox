package mqtt

import (
	"github.com/IvoryRaptor/iotbox/common"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

// AMqtt 上报模块
type AMqtt struct {
	common.AModule
	// mqtt 服务器地址
	address string
	// client id
	clientID string
	// 用户名
	user string
	// 密码
	password string
	// 上报主题
	topic string
	// 服务质量
	qos byte
	// mqtt client
	client MQTT.Client
}

// GetName 获取模块名
func (m *AMqtt) GetName() string {
	return "mqtt"
}

// Config 配置模块
func (m *AMqtt) Config(_ common.IKernel, config map[string]interface{}) error {
	if _, ok := config["address"]; ok {
		m.address = config["address"].(string)
	}

	if _, ok := config["clientID"]; ok {
		m.clientID = config["clientID"].(string)
	}

	if _, ok := config["user"]; ok {
		m.user = config["user"].(string)
	}

	if _, ok := config["password"]; ok {
		m.password = config["password"].(string)
	}
	if _, ok := config["topic"]; ok {
		m.topic = config["topic"].(string)
	}
	if _, ok := config["qos"]; ok {
		m.qos = byte(config["qos"].(int))
	} else {
		m.qos = 0
	}

	log.Printf("[%s]==> Config %#v\n", m.GetName(), m)
	if err := m.createConnect(); err != nil {
		log.Fatalf("[%s]===> %s", m.GetName(), err)
	}
	return nil
}

// Send 发送数据
func (m *AMqtt) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	log.Printf("[%s]==> Send %s\n", m.GetName(), string(packet["value"].([]byte)))
	go func() {
		if !m.client.IsConnected() {
			if err := m.createConnect(); err != nil {
				log.Fatalf("[%s]===> %s", m.GetName(), err)
				m.Response <- nil
				return
			}
		}
		// text := `
		// {
		//   "DO01": true,
		//   "AO01": -121,
		//   "AO02": 12341,
		//   "AI03": 456,
		//   "AI04": 3.14,
		//   "AI05": 3.14544
		// }
		// `
		token := m.client.Publish(m.topic, m.qos, false, packet["value"])
		token.Wait()
		if err := token.Error(); err != nil {
			log.Fatalf("[%s]===> %s", m.GetName(), err)
		}
		m.Response <- nil
		return
	}()
	return m.Response
}

// Create 创建mqtt 上报对象
func Create() *AMqtt {
	log.Println("Create upload mqtt module")
	return &AMqtt{}
}

func (m *AMqtt) defaultMessageHandler(client MQTT.Client, msg MQTT.Message) {
	log.Printf("[%s]===> topic: %s\n", m.GetName(), msg.Topic())
	log.Printf("[%s]===> msg: %s\n", m.GetName(), msg.Payload())
}

// createConnect 创建mqtt连接
func (m *AMqtt) createConnect() error {
	opts := MQTT.NewClientOptions().AddBroker(m.address)
	opts.SetClientID(m.clientID)
	opts.SetUsername(m.user)
	opts.SetPassword(m.password)
	opts.SetDefaultPublishHandler(m.defaultMessageHandler)

	//create and start a client
	m.client = MQTT.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
