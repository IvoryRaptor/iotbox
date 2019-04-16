package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/common"
	"sync"
	"time"
)

type MockPassivePort struct {
	name string
}

func (m *MockPassivePort) Open() error {
	return nil
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
	msg   common.Message
	mutex *sync.Mutex
}

func (m *MockActivePort) Open() error {
	m.mutex = new(sync.Mutex)
	return nil
}

func (m *MockActivePort) Write(message common.Message) error {
	return nil
}

func (m *MockActivePort) Read(wait time.Duration) (msg common.Message, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	result := m.msg
	if result == nil {
		time.Sleep(wait)
	} else {
		println("2")
	}
	m.msg = nil
	return result, nil
}

func (m *MockActivePort) Close() error {
	return nil
}

func (m *MockActivePort) SetMessage(msg common.Message) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	println(1)
	m.msg = msg
}

func main() {
	active := &MockActivePort{}
	active.Open()

	common.CreateActivePort("net", active, 1*time.Second, func(message common.Message) {
		common.AddTask("com1", &ReadTask{
			Owner:   "net",
			Message: common.Message{"name": "g"},
			Wait:    1 * time.Second,
		})
	})
	common.CreatePassivePort("com1", &MockPassivePort{})

	//common.AddTask("com1", &TestTask{})
	//common.AddTask("com1", &TestTask{})
	//
	//common.AddTask("com1", &ArrayTask{
	//	Messages: []common.Message{
	//		{"name": "d"},
	//		{"name": "e"},
	//		{"name": "f"},
	//	},
	//	Wait: 1 * time.Second,
	//})
	for {
		console.ReadLine()
		active.SetMessage(common.Message{})
	}
}
