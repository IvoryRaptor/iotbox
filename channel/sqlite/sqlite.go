package sqlite

import "github.com/IvoryRaptor/iotbox/common"

type Sqlite struct {
}

func (m *Sqlite) Config(ch chan common.ITask, config map[string]interface{}) error {

}

func (m *Sqlite) Send(packet common.Packet) chan common.Packet {

}
