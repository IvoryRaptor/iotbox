package downsidemock

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"math/rand"
	"time"
)

type Mock struct {
	common.AModule
}

func (m *Mock) GetName() string {
	return "mock"
}

func (m *Mock) Config(_ common.IKernel, config map[string]interface{}) error {
	return nil
}

func (m *Mock) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	fmt.Printf("[downsidemock] Send Packet %s\n", packet["address"])
	b := rand.Intn(3)
	if b > 1 {
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Printf("[downsidemock] Receive Packet\n")
			m.Response <- common.Packet{
				"value": rand.Intn(100),
			}
		}()
	}
	return m.Response
}

func Create() *Mock {
	return &Mock{}
}
