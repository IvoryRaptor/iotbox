package task

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/mocktask"
)

func CreateTask(kernel common.IKernel, config map[string]interface{}) error {
	taskType := config["type"]
	switch taskType {
	case "mock":
		mock := &mocktask.MockTask{}
		if err := mock.Config(kernel, config); err != nil {
			return err
		}
		fmt.Printf("Add Task [%s] %s\n", taskType, config["cron"])
		kernel.JoinTask(config["cron"].(string), mock)
		return nil
	}
	return nil
}
