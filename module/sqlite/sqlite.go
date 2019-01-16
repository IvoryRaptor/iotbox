package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	db   *sql.DB
	read chan common.Packet
}

func (m *Sqlite) Config(ch chan common.ITask, config map[string]interface{}) error {
	var err error
	m.read = make(chan common.Packet)
	m.db, err = sql.Open("sqlite3", config["filename"].(string))
	for _, sqlText := range config["init"].([]interface{}) {
		if _, err := m.db.Exec(sqlText.(string)); err != nil {
			return err
		}
	}
	go func() {
		for {
			task := <-ch
			task.Work(m)
		}
	}()
	return err
}

func (m *Sqlite) Send(packet common.Packet) chan common.Packet {
	go func() {
		fmt.Printf("SQLITE EXEC [%s]\n", packet["sql"])
		if result, err := m.db.Exec(packet["sql"].(string)); err != nil {
			m.read <- common.Packet{
				"error": err,
			}
		} else {
			m.read <- common.Packet{
				"result": result,
			}
		}
	}()
	return m.read
}
