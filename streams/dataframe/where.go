package dataframe

import (
	"github.com/IvoryRaptor/iotbox/akka"
)

type Where struct {
	DataFrame
	rows  []Row
	count int
	work  func(row Row) bool
}

func (w *Where) Receive(context akka.Context) {
	switch msg := context.Message().(type) {
	case *akka.Started:
		//w.init(context)
	case *InsertRow:
		row := msg.Row
		if w.work(row) {
			//w.insert(msg)
		}
	case *RemoveRow:
		//w.remove(msg)
	}
}
