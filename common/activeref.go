package common

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/akka"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type ActiveRef struct {
	Port Port
	Wait time.Duration
	Work func(message Message)
}

type Idle struct {
}

func (a *ActiveRef) Receive(context akka.Context) {
	switch task := context.Message().(type) {
	case *akka.Started:
		println(context.Self().Id)
		data, err := ioutil.ReadFile(fmt.Sprintf("./config/port/%s", context.Self().Id))
		if err != nil {
			data = []byte{}
		}
		var config map[string]interface{}
		if err := yaml.Unmarshal(data, &config); err != nil {

		}
		a.Port.Open(config)
	case *TaskRef:
		request := task.GetRequest()
		for request != nil {
			retry := 0
			for ; retry < request.Retry; retry++ {
				a.Port.Write(request.Body)
				if result, err := a.Port.Read(request.Wait); err == nil {
					task.Receive(&Response{
						State: ResponseTimeout,
						Body:  result,
					})
					break
				}
			}
			if retry >= request.Retry {
				task.Receive(&Response{
					State: ResponseResult,
					Body:  nil,
				})
			}
			request = task.GetRequest()
		}
	case *Idle:
		if result, err := a.Port.Read(a.Wait); err != nil {

		} else if result != nil {
			println("work")
			a.Work(result)
		}
		context.Tell(context.Self(), context.Message())
	}
}
