package kernel

import (
	"errors"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/demo"
	"github.com/IvoryRaptor/iotbox/task/owner"
	"github.com/IvoryRaptor/iotbox/task/report"
	"github.com/IvoryRaptor/iotbox/task/sql"
)

func (k *Kernel) CreateTask(config map[interface{}]interface{}) (common.ITask, error) {
	taskType := config["type"].(string)
	var result common.ITask
	switch taskType {
	case "mock":
		result = demo.CreateDemo()
	case "sql":
		result = sql.CreateSql()
	case "report":
		result = report.CreateReport()
	case "owner":
		result = &owner.Owner{}
	}
	if result == nil {
		return nil, errors.New(fmt.Sprintf("Unknown Task Type [%s]", taskType))
	}
	if err := result.Config(k, config); err != nil {
		return nil, err
	}
	return result, nil
}
