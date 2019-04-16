package port

import (
	"github.com/IvoryRaptor/iotbox/common"
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
