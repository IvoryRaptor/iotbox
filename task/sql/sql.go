package sql

import (
	"bytes"
	"github.com/IvoryRaptor/iotbox/common"
	"text/template"
)

type Sql struct {
	tpl    *template.Template
	packet common.Packet
	target string
	kernel common.IKernel
}

func (s *Sql) Config(kernel common.IKernel, config map[interface{}]interface{}) error {
	var err error
	s.tpl, err = template.New("").Parse(config["sql"].(string))
	s.target = config["target"].(string)
	s.kernel = kernel
	return err
}

func (s *Sql) SetPacket(packet common.Packet) common.IHandlerTask {
	s.packet = packet
	return s
}

func (s *Sql) Run() {
	s.kernel.GetModule(s.target) <- s
}

func (s *Sql) Clone() common.IHandlerTask {
	return &Sql{
		tpl:    s.tpl,
		target: s.target,
		kernel: s.kernel,
	}
}

func (s *Sql) Work(channel common.IModule) {
	buf := new(bytes.Buffer)
	if err := s.tpl.Execute(buf, s.packet); err != nil {
		return
	}
	ch := channel.Send(s, common.Packet{
		"sql": buf.String(),
	})
	<-ch
}
