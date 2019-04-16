package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type MockPassivePort struct {
	name string
}

func (m *MockPassivePort) Write(message common.Message) error {
	m.name = message["name"].(string)
	println("write:" + m.name)
	return nil
}

func (m *MockPassivePort) Read(wait time.Duration) (msg common.Message, err error) {
	println("read:" + m.name)
	return common.Message{
		"name":  m.name,
		"value": "1",
	}, nil
}

func (m *MockPassivePort) Close() error {
	return nil
}

type MockActivePort struct {
	msg common.Message
}

func (m *MockActivePort) Write(message common.Message) error {
	return nil
}

func (m *MockActivePort) Read(wait time.Duration) (msg common.Message, err error) {
	result := m.msg
	m.msg = nil
	return result, nil
}

func (m *MockActivePort) Close() error {
	return nil
}

func (m *MockActivePort) SetMessage(msg common.Message) {
	m.msg = msg
}

func main() {
	common.CreateActivePort("net", &MockActivePort{})
	common.CreatePassivePort("com1", &MockPassivePort{})

	common.AddTask("com1", &TestTask{})
	common.AddTask("com1", &TestTask{})

	common.AddTask("com1", &ArrayTask{
		Messages: []common.Message{
			{"name": "d"},
			{"name": "e"},
			{"name": "f"},
		},
		Wait: 1 * time.Second,
	})

	common.AddTask("com1", &ReadTask{
		Owner:   "net",
		Message: common.Message{"name": "g"},
		Wait:    1 * time.Second,
	})
	console.ReadLine()
}
