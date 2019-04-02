package sqlite

import (
	"database/sql"
	"errors"
	"github.com/IvoryRaptor/iotbox/test/akka"
)

type SqliteActor struct {
	akka.Actor
	db       *sql.DB
	filename string
	execSql  []string
}

func (actor *SqliteActor) PreStart() error {
	var err error
	actor.db, err = sql.Open("sqlite3", actor.filename)
	for _, sqlText := range actor.execSql {
		if _, err := actor.db.Exec(sqlText); err != nil {
			return err
		}
	}
	return err
}

func (actor *SqliteActor) Config(config map[string]interface{}) error {
	actor.filename = config["filename"].(string)
	actor.execSql = make([]string, len(config["init"].([]interface{})))
	for i, sqlText := range config["init"].([]interface{}) {
		actor.execSql[i] = sqlText.(string)
	}
	return nil
}

func (actor *SqliteActor) Receive(sender akka.IActor, message akka.Message) error {
	if exec, hive := message.GetString("exec"); hive {
		if _, err := actor.db.Exec(exec); err != nil {
			return err
		} else {
			return nil
		}
	} else if query, hive := message.GetString("query"); hive {
		if rows, err := actor.db.Query(query); err != nil {
			return err
		} else {
			for rows.Next() {
				row := map[string]interface{}{}
				columns, _ := rows.Columns()
				for _, column := range columns {
					row[column] = nil
				}
				rows.Scan()
				//rows.Scan()
			}
		}
	}
	return errors.New("")
}
