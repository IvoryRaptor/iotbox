package kernel

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task/demo"
	"github.com/IvoryRaptor/iotbox/task/owner"
	"github.com/IvoryRaptor/iotbox/task/report"
	"github.com/IvoryRaptor/iotbox/task/sql"
	"github.com/IvoryRaptor/iotbox/task/qa"
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
	case "QA":
		result = qa.Create()
	}
	if result == nil {
		return nil, fmt.Errorf(fmt.Sprintf("Unknown Task Type [%s]", taskType))
	}
	if err := result.Config(k, config); err != nil {
		return nil, err
	}
	return result, nil
}
