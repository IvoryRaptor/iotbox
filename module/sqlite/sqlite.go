package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	common.AModule
	db *sql.DB
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

func (m *Sqlite) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	go func() {
		fmt.Printf("SQLITE EXEC [%s]\n", packet["sql"])
		if result, err := m.db.Exec(packet["sql"].(string)); err != nil {
			m.Response <- common.Packet{
				"error": err,
			}
		} else {
			m.Response <- common.Packet{
				"result": result,
			}
		}
	}()
	return m.Response
}
