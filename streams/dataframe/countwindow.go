package dataframe

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/akka"
)

type CountWindow struct {
	DataFrame
	rows        []Row
	index       int
	count       int
	WindowCount int
}

func (t *CountWindow) GetRows() []Row {
	if t.count >= t.WindowCount {
		return t.rows
	} else {
		return t.rows[0:t.index]
	}
}

func (t *CountWindow) init(context akka.Context) {
	t.index = 0
	t.count = 0
	t.rows = make([]Row, t.WindowCount)
	t.Start(context)
}

func (t *CountWindow) insert(msg *InsertRow) {
	if t.count < t.WindowCount {
		t.count++
	} else {
		t.TellNext(&RemoveRow{
			Row: t.rows[t.index],
		})
	}
	t.rows[t.index] = msg.Row
	t.index = t.index + 1
	if t.index >= t.WindowCount {
		t.index = 0
	}
	t.TellNext(msg)
}

func (t *CountWindow) remove(msg *RemoveRow) {
	move := false
	row := msg.Row
	for index, r := range t.rows {
		if move {
			t.rows[index-1] = t.rows[index]
		} else {
			if r == row {
				move = true
				t.rows[index] = nil
				t.count = t.count - 1
				t.index = (t.index + t.WindowCount - 1) % t.WindowCount
			}
		}
	}
}

func (t *CountWindow) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		t.init(context)
	case *InsertRow:
		t.insert(msg)
	case *RemoveRow:
		t.remove(msg)
	}
}

func (t *CountWindow) Show() {
	for _, row := range t.GetRows() {
		for _, col := range row.GetValues() {
			fmt.Printf("%d\t", col.(int))
		}
		fmt.Println()
	}
}

func (t *CountWindow) Where(sql string) *Where {
	result := &Where{
		DataFrame: DataFrame{
			Schema: t.Schema,
		},
		work: func(row Row) bool {
			field1, value1 := row.GetValue("A")
			_, value2 := row.GetValue("B")
			return field1.DataType.Equal(value1, value2)
		},
	}
	t.Append(result)
	return result
}
