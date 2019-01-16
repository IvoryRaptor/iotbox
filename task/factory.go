package task

import (
	"errors"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/mock"
	"github.com/IvoryRaptor/iotbox/task/sql"
)

func CreateTask(kernel common.IKernel, config map[interface{}]interface{}) (common.ITask, error) {
	taskType := config["type"].(string)
	var result common.ITask
	switch taskType {
	case "mock":
		result = &mock.Mock{}
	case "sql":
		result = &sql.Sql{}
	}
	if result == nil {
		return nil, errors.New(fmt.Sprintf("Unknown Task Type [%s]", taskType))
	}
	if err := result.Config(kernel, config); err != nil {
		return nil, err
	}
	return result, nil
}
