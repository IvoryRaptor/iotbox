package sql

import (
	"bytes"
	"github.com/IvoryRaptor/iotbox/common"
	"text/template"
)

type Sql struct {
	receiver chan common.ITask
	tpl      *template.Template
	packet   common.Packet
}

func (s *Sql) Config(kernel common.IKernel, config map[interface{}]interface{}) error {
	var err error
	s.tpl, err = template.New("").Parse(config["sql"].(string))
	s.receiver = kernel.GetModule(config["receiver"].(string))
	return err
}

func (s *Sql) Run() {
	s.receiver <- s
}

func (s *Sql) SetPacket(packet common.Packet) common.IHandlerTask {
	s.packet = packet
	return s
}

func (s *Sql) Clone() common.IHandlerTask {
	return &Sql{
		receiver: s.receiver,
		tpl:      s.tpl,
	}
}

func (s *Sql) Work(channel common.IModule) {
	buf := new(bytes.Buffer)
	if err := s.tpl.Execute(buf, s.packet); err != nil {
		return
	}
	ch := channel.Send(common.Packet{
		"sql": buf.String(),
	})
	<-ch
}
