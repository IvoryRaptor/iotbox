package mqtt

import (
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/protocol"
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
	// 服务质量
	qos byte
	// 协议
	protocolType string
	// 订阅主题
	subscribe []string
	// 推送主题
	publish []string
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
	m.subscribe = make([]string, 0)
	if _, ok := config["subscribe"]; ok && config["subscribe"] != nil {
		for _, v := range config["subscribe"].([]interface{}) {
			m.subscribe = append(m.subscribe, v.(string))
		}

	}
	m.publish = make([]string, 0)
	if _, ok := config["publish"]; ok && config["publish"] != nil {
		for _, v := range config["publish"].([]interface{}) {
			m.publish = append(m.publish, v.(string))
		}
	}
	if _, ok := config["protocolType"]; ok {
		m.protocolType = config["protocolType"].(string)
	}
	if _, ok := config["qos"]; ok {
		m.qos = byte(config["qos"].(int))
	} else {
		m.qos = 0
	}
	log.Printf("[%s]==> Config %#v\n", m.GetName(), m)
	if err := m.createConnect(); err != nil {
		log.Printf("[%s]======> %s", m.GetName(), err)
	}
	return nil
}

// Send 发送数据
func (m *AMqtt) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	log.Printf("[%s]==> Send", m.GetName())
	go func() {
		var p protocol.IProtocol
		var err error
		var sendBuf []byte
		if !m.client.IsConnected() {
			if err = m.createConnect(); err != nil {
				goto breakout
			}
		}
		p, err = protocol.CreateProtocol("rootcloud", nil)
		if err != nil {
			goto breakout
		}
		p.Config(packet)
		sendBuf, err = p.Encode(packet)
		if err != nil {
			goto breakout
		}
		for _, item := range m.publish {
			token := m.client.Publish(item, m.qos, false, sendBuf)
			token.Wait()
			if err := token.Error(); err != nil {
				log.Printf("[%s]===>top[%s] %s", m.GetName(), item, err)
			}
		}
		m.Response <- nil
		return
	breakout:
		log.Printf("[%s]====> %s", m.GetName(), err)
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
func (m *AMqtt) createConnect() (err error) {
	opts := MQTT.NewClientOptions().AddBroker(m.address)
	opts.SetClientID(m.clientID)
	opts.SetUsername(m.user)
	opts.SetPassword(m.password)
	opts.SetDefaultPublishHandler(m.defaultMessageHandler)

	//create and start a client
	m.client = MQTT.NewClient(opts)
	token := m.client.Connect()
	token.Wait()
	err = token.Error()
	return err
}

// // createProtocol 创建协议
// func (m *AMqtt) createProtocol() (common.IProtocol, error) {
// 	var res common.IProtocol
// 	switch strings.ToLower(m.protocolType) {
// 	case "rootcloud":
// 		res = rootcloud.Create()
// 	default:
// 		return nil, fmt.Errorf("protocolType not support[%s]", m.protocolType)
// 	}
// 	return res, nil
// }
>>>>>>> origin/master
