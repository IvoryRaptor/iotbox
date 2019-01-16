package mock

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"math/rand"
	"time"
)

type Mock struct {
	common.AModule
}

func (m *Mock) Config(ch chan common.ITask, config map[string]interface{}) error {
	m.Start(ch, m)
	return nil
}

func (m *Mock) Send(packet common.Packet) chan common.Packet {
	fmt.Printf("[mock] Send Packet\n")
	b := rand.Intn(3)
	if b > 1 {
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Printf("[mock] Receive Packet\n")
			m.Response <- common.Packet{
				"value": rand.Intn(100),
			}
		}()
	}
	return m.Response
}
