package port

import (
	"github.com/IvoryRaptor/iotbox/common"
	"sync"
	"time"
)

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
