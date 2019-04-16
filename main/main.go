package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/port"
	"github.com/IvoryRaptor/iotbox/task"
	"time"
)

func main() {
	active := &port.MockActivePort{}

	common.CreateActivePort("net", active, 1*time.Second, func(message common.Message) {
		common.AddTask("com1", &task.ReadTask{
			Owner:   "net",
			Message: common.Message{"name": "g"},
			Wait:    1 * time.Second,
		})
	})

	common.CreatePassivePort("com1", &port.MockPassivePort{})

	//common.AddTask("com1", &TestTask{})
	//common.AddTask("com1", &TestTask{})
	//
	//common.AddTask("com1", &ArrayTask{
	//	Messages: []common.Message{
	//		{"name": "d"},
	//		{"name": "e"},
	//		{"name": "f"},
	//	},
	//	Wait: 1 * time.Second,
	//})
	for {
		console.ReadLine()
		active.SetMessage(common.Message{})
	}
}
