package sql

import (
	"bytes"
	"github.com/IvoryRaptor/iotbox/common"
	"text/template"
	"time"
)

type Sql struct {
	common.ATask
	tpl        *template.Template
	sqlCommand string
}

func (s *Sql) SetPacket(packet common.Packet) (common.ICloneTask, error) {
	buf := new(bytes.Buffer)
	var err error
	if err = s.tpl.Execute(buf, packet); err == nil {
		s.sqlCommand = buf.String()
	}
	return s, err
}

func (s *Sql) Clone() common.ICloneTask {
	result := &Sql{
		tpl:   s.tpl,
		ATask: s.ATask,
	}
	result.SetCurrentWork(result.SqlWork)
	return result
}

func (s *Sql) SqlConfig(kernel common.IKernel, config map[interface{}]interface{}) error {
	var err error
	s.SetCurrentWork(s.SqlWork)
	s.tpl, err = template.New("").Parse(config["sql"].(string))
	return err
}

func (s *Sql) SqlWork(module common.IModule) (common.WorkState, error) {
	ch := module.Send(s, common.Packet{
		"sql": s.sqlCommand,
	})
	//消费掉回应消息
	module.Read(ch, time.Second*3)
	return common.Complete, nil
}

func CreateSql() *Sql {
	result := &Sql{}
	result.SetOtherConfig(result.SqlConfig)
	return result
}
