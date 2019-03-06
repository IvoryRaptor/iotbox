package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Sqlite struct {
	common.AModule
	db *sql.DB
}

func (m *Sqlite) GetName() string {
	return "sqlite"
}

func (m *Sqlite) Config(_ common.IKernel, config map[string]interface{}) error {
	var err error
	m.db, err = sql.Open("sqlite3", config["filename"].(string))
	for _, sqlText := range config["init"].([]interface{}) {
		if _, err := m.db.Exec(sqlText.(string)); err != nil {
			return err
		}
	}
	return err
}

// Send sqlite 执行体
func (m *Sqlite) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	go func() {
		var err error
		var value []string
		// 判断类型是sql
		if _, ok := packet["type"]; !ok {
			err = fmt.Errorf("not find type")
			goto breakout
		}
		if _, ok := packet["value"]; !ok {
			err = fmt.Errorf("not find value")
			goto breakout
		}
		value = packet["value"].([]string)
		for _, item := range value {
			log.Printf("[sql]==> %s", item)
			if _, err = m.db.Exec(item); err != nil {
				log.Printf("[sql]==> err[%s]", err)
			}
		}
		m.Response <- common.Packet{
			"type":   "sql",
			"status": "ok",
			"desc":   nil,
			"value":  nil,
		}
		return
	breakout:
		m.Response <- common.Packet{
			"type":   "sql",
			"status": "error",
			"desc":   err,
			"value":  nil,
		}
	}()
	return m.Response
}

func CreateSqlite() *Sqlite {
	return &Sqlite{}
}
