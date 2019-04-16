package common

import "github.com/IvoryRaptor/iotbox/akka"

type ActiveRef struct {
	Port Port
}

func (a *ActiveRef) Receive(context akka.Context) {
	switch task := context.Message().(type) {
	case *akka.Started:
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
	}
}
