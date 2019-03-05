package sql

import (
	"bytes"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"text/template"
	"time"
)

type Sql struct {
	common.ATask
	tpl        *template.Template
	sqlCommand []string
}

// SetPacket 设置包
func (s *Sql) SetPacket(packet common.Packet) (common.ICloneTask, error) {
	var value []map[string]interface{}
	if _, ok := packet["type"]; !ok {
		return nil, fmt.Errorf("not find type")
		// cType = packet["type"].(string)
	}

	if _, ok := packet["value"]; ok {
		value = packet["value"].([]map[string]interface{})
	} else {
		return nil, fmt.Errorf("not find value")
	}

	s.sqlCommand = make([]string, 0)
	for _, item := range value {
		buf := new(bytes.Buffer)
		if err := s.tpl.Execute(buf, item); err == nil {
			s.sqlCommand = append(s.sqlCommand, buf.String())
		}
	}
	return s, nil
}

func (s *Sql) Clone() common.ICloneTask {
	result := &Sql{
		tpl:   s.tpl,
		ATask: s.ATask,
	}
	result.SetCurrentWork(result.SqlWork)
	return result
}

func (s *Sql) SqlConfig(kernel common.IKernel, config map[string]interface{}) error {
	var err error
	s.SetCurrentWork(s.SqlWork)
	s.tpl, err = template.New("").Parse(config["sql"].(string))
	return err
}

func (s *Sql) SqlWork(module common.IModule) (common.WorkState, error) {
	ch := module.Send(s, common.Packet{
		"type":  "sql",
		"value": s.sqlCommand,
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
