package common

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/akka"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type PassiveRef struct {
	Port Port
}

func (p *PassiveRef) Receive(context akka.Context) {
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
		p.Port.Open(config)
	case *TaskRef:
		request := task.GetRequest()
		for request != nil {
			retry := 0
			for ; retry < request.Retry; retry++ {
				p.Port.Write(request.Body)
				if result, err := p.Port.Read(request.Wait); err == nil {
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
	}
}
