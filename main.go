package main

import (
	"fmt"
	"github.com/AsynkronIT/goconsole"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/streams"
)

func main() {
	source := streams.Source{}

	source.Foreach(func(msg interface{}) {
		m := msg.(common.Message)
		println("id=" + m["id"].(string))
	})

	source.Map(func(msg interface{}) interface{} {
		m := msg.(common.Message)
		m["result"] = m["a"].(int) + m["b"].(int)
		return msg
	}).Filter(func(msg interface{}) bool {
		m := msg.(common.Message)
		return m["result"].(int)%2 == 0
	}).Foreach(func(msg interface{}) {
		m := msg.(common.Message)
		println(m["result"].(int))
	})

	source.
		Window(5).
		Map(func(msg []interface{}) interface{} {
			sum := 0
			for _, m := range msg {
				sum += m.(common.Message)["a"].(int)
			}
			return common.Message{"sum": sum}
		}).Foreach(func(msg interface{}) {
		m := msg.(common.Message)
		println(m["sum"].(int))
	})

	source.Start()

	for i := 0; i < 10; i++ {
		for j := 1; j < 10; j++ {
			source.Write(common.Message{
				"id": fmt.Sprintf("%d.%d", i, j),
				"a":  i,
				"b":  j,
			})
		}
	}

	for {
		console.ReadLine()
	}
}
