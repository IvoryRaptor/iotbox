package kernel

import (
	"errors"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/module/corn"
	"github.com/IvoryRaptor/iotbox/module/mock"
	"github.com/IvoryRaptor/iotbox/module/sqlite"
)

func (k *Kernel) CreateModule(config map[string]interface{}) (common.IModule, error) {
	channelType := config["type"].(string)
	var result common.IModule
	switch channelType {
	case "mock":
		result = &mock.Mock{}
	case "sqlite":
		result = &sqlite.Sqlite{}
	case "corn":
		result = &corn.Corn{}
	}
	if result == nil {
		return nil, errors.New(fmt.Sprintf("Unknown Module Type [%s]", channelType))
	}
	if err := result.Config(k, config); err != nil {
		return nil, err
	}
	return result, nil
}
