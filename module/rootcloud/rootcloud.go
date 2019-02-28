package rootcloud

import (
	"github.com/IvoryRaptor/iotbox/common"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

// RootCloud 树根云平台接入
// http://www.rootcloud.com/
type RootCloud struct {
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
func (rc *RootCloud) GetName() string {
	return "RootCloud"
}

// Config 配置模块
func (rc *RootCloud) Config(_ common.IKernel, config map[string]interface{}) error {
	if _, ok := config["address"]; ok {
		rc.address = config["address"].(string)
	}

	if _, ok := config["clientID"]; ok {
		rc.clientID = config["clientID"].(string)
	}

	if _, ok := config["user"]; ok {
		rc.user = config["user"].(string)
	}

	if _, ok := config["password"]; ok {
		rc.password = config["password"].(string)
	}
	if _, ok := config["topic"]; ok {
		rc.topic = config["topic"].(string)
	}
	if _, ok := config["qos"]; ok {
		rc.qos = byte(config["qos"].(int))
	} else {
		rc.qos = 0
	}

	log.Printf("[%s]==> Config %#v\n", rc.GetName(), rc)
	if err := rc.createConnect(); err != nil {
		log.Fatalf("[%s]===> %s", rc.GetName(), err)
	}
	return nil
}

// Send 发送数据
func (rc *RootCloud) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	log.Printf("[%s]==> Send\n", rc.GetName())
	go func() {
		if !rc.client.IsConnected() {
			rc.createConnect()
		}
		text := `
		{
		  "DO01": true,
		  "AO01": -121,
		  "AO02": 12341,
		  "AI03": 456,
		  "AI04": 3.14,
		  "AI05": 3.14544
		}
		`
		token := rc.client.Publish(rc.topic, rc.qos, false, text)
		token.Wait()
		if err := token.Error(); err != nil {
			log.Fatalf("[%s]===> %s", rc.GetName(), err)
		}
		rc.Response <- nil
		return
	}()
	return rc.Response
}

// Create 创建树根云上报对象
func Create() *RootCloud {
	log.Println("Create RootCloud module")
	return &RootCloud{}
}

func (rc *RootCloud) defaultMessageHandler(client MQTT.Client, msg MQTT.Message) {
	log.Printf("[%s]===> topic: %s\n", rc.GetName(), msg.Topic())
	log.Printf("[%s]===> msg: %s\n", rc.GetName(), msg.Payload())
}

// createConnect 创建mqtt连接
func (rc *RootCloud) createConnect() error {
	opts := MQTT.NewClientOptions().AddBroker(rc.address)
	opts.SetClientID(rc.clientID)
	opts.SetUsername(rc.user)
	opts.SetPassword(rc.password)
	opts.SetDefaultPublishHandler(rc.defaultMessageHandler)

	//create and start a client
	rc.client = MQTT.NewClient(opts)
	if token := rc.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
