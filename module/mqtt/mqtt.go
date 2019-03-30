package mqtt

import (
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	common.AModule
	client mqtt.Client
}

func (m *MQTT) GetName() string {
	return "mqtt"
}

func (m *MQTT) Config(_ common.IKernel, config map[string]interface{}) error {
	opts := mqtt.NewClientOptions().AddBroker("tcp://iot.eclipse.org:1883").SetClientID("gotrivial")
	m.client = mqtt.NewClient(opts)

	return nil
}

func (m *MQTT) Send(_ common.ITask, packet common.Packet) chan common.Packet {

	return m.Response
}
