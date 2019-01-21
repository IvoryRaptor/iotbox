package sql

import (
	"bytes"
	"github.com/IvoryRaptor/iotbox/common"
	"text/template"
)

type Sql struct {
	common.ATask
	tpl    *template.Template
	packet common.Packet
}

func (s *Sql) SetPacket(packet common.Packet) common.IHandlerTask {
	s.packet = packet
	return s
}

func (s *Sql) Clone() common.IHandlerTask {
	result := &Sql{
		tpl:   s.tpl,
		ATask: s.ATask,
	}
	return InitSql(result)
}

func (s *Sql) SqlConfig(kernel common.IKernel, config map[interface{}]interface{}) error {
	var err error
	s.tpl, err = template.New("").Parse(config["sql"].(string))
	return err
}

func (s *Sql) SqlWork(channel common.IModule) (common.WorkState, error) {
	buf := new(bytes.Buffer)
	if err := s.tpl.Execute(buf, s.packet); err != nil {
		return common.Failed, err
	}
	ch := channel.Send(s, common.Packet{
		"sql": buf.String(),
	})
	<-ch
	return common.Complete, nil
}

func InitSql(sql *Sql) *Sql {
	sql.SetCurrentWork(sql.SqlWork).SetOtherConfig(sql.SqlConfig)
	return sql
}

func Create() *Sql {
	result := &Sql{}
	return InitSql(result)
}
