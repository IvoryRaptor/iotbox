package upsidemock

import (
	"github.com/IvoryRaptor/iotbox/common"
)

type Upside struct {
	common.AModule
}

func (m *Upside) Config(_ common.IKernel, config map[string]interface{}) error {
	return nil
}

func (m *Upside) Send(packet common.Packet) chan common.Packet {

	return m.Response
}
