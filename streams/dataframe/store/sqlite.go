package store

import (
	"database/sql"
	"fmt"
	"github.com/IvoryRaptor/iotbox/streams/dataframe"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

type Sqlite struct {
	db *sql.DB
}

func (f *Sqlite) CreateTable(schema *dataframe.Schema) string {
	var sqlText strings.Builder
	fmt.Fprint(&sqlText, "CREATE TABLE IF NOT EXISTS DAT_MOCK(ID INTEGER PRIMARY KEY AUTOINCREMENT,VALUE INTEGER NOT NULL);")
	fmt.Fprint(&sqlText, schema.Name)
	for _, field := range schema.Fields {
		fmt.Fprint(&sqlText, field.Name)
		switch field.DataType.(type) {
		case *dataframe.BooleanType:
			fmt.Fprint(&sqlText, " BOOL ")
			if field.Nullable {
				fmt.Fprint(&sqlText, " NOT NULL ")
			}
		case *dataframe.IntegerType:

		}
	}
	fmt.Fprint(&sqlText, schema.Name)
	return sqlText.String()
}

func (f *Sqlite) Insert(row *dataframe.Row) {
	var err error
	schema := row.Schema

	var sqlText strings.Builder
	fmt.Fprint(&sqlText, "INSERT INOT ")
	fmt.Fprint(&sqlText, schema.Name+"(")

	for index, field := range row.Schema.Fields {
		fmt.Fprint(&sqlText, field.Name)
		if index != len(row.Schema.Fields)-1 {
			fmt.Fprint(&sqlText, ",")
		}
	}
	fmt.Fprint(&sqlText, schema.Name+") VALUES (")

	for index, field := range row.Schema.Fields {
		fmt.Fprint(&sqlText, field.DataType.ValueToString(row.Values[index]))
		if index != len(row.Schema.Fields)-1 {
			fmt.Fprint(&sqlText, ",")
		}
	}

	fmt.Fprint(&sqlText, ")")
	if _, err = f.db.Exec(sqlText.String()); err != nil {
		log.Printf("[sqlText]==> err[%s]", err)
	}
}
