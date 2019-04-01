package downsidemock

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"math/rand"
	"time"
)

type Mock struct {
	common.AModule
	failureRate int
	wait        time.Duration
}

func (m *Mock) GetName() string {
	return "mock"
}

func (m *Mock) Config(_ common.IKernel, config map[string]interface{}) error {
	m.failureRate = config["failure"].(int) - 1
	m.wait = time.Duration(config["wait"].(int)) * time.Second
	return nil
}

func (m *Mock) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	fmt.Printf("[downsidemock] Tell Packet %s\n", packet["address"])
	b := rand.Intn(100)
	if b > m.failureRate {
		go func() {
			time.Sleep(m.wait)
			fmt.Printf("[downsidemock] Receive Packet\n")
			m.Response <- common.Packet{
				"value": rand.Intn(100),
			}
		}()
	}
	return m.Response
}

func CreateMock() *Mock {
	return &Mock{}
}
