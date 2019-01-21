package kernel

import (
	"errors"
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/module/corn"
	"github.com/IvoryRaptor/iotbox/module/downsidemock"
	"github.com/IvoryRaptor/iotbox/module/sqlite"
	"github.com/IvoryRaptor/iotbox/module/upsidemock"
)

func (k *Kernel) CreateModule(config map[string]interface{}) (common.IModule, error) {
	channelType := config["type"].(string)
	var result common.IModule
	switch channelType {
	case "downsidemock":
		result = downsidemock.Create()
	case "sqlite":
		result = sqlite.Create()
	case "corn":
		result = corn.Create()
	case "upsidemock":
		result = upsidemock.Create()
	}
	if result == nil {
		return nil, errors.New(fmt.Sprintf("Unknown Module Type [%s]", channelType))
	}
	if err := result.Config(k, config); err != nil {
		return nil, err
	}
	return result, nil
}
