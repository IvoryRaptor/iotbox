package kernel

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/module/corn"
	"github.com/IvoryRaptor/iotbox/module/downsidemock"
	"github.com/IvoryRaptor/iotbox/module/sqlite"
	"github.com/IvoryRaptor/iotbox/module/upsidemock"
	"github.com/IvoryRaptor/iotbox/module/modbus"
	"github.com/IvoryRaptor/iotbox/module/mqtt"
)

func (k *Kernel) CreateModule(config map[string]interface{}) (common.IModule, error) {
	channelType := config["type"].(string)
	var result common.IModule
	switch channelType {
	case "downsidemock":
		result = downsidemock.CreateMock()
	case "sqlite":
		result = sqlite.CreateSqlite()
	case "corn":
		result = corn.CreateCore()
	case "upsidemock":
		result = upsidemock.CreateUpside()
	case "modbus":
		result = modbus.Create()
	case "mqtt":
		result = mqtt.Create()
	}
	if result == nil {
		return nil, fmt.Errorf(fmt.Sprintf("Unknown Module Type [%s]", channelType))
	}
	if err := result.Config(k, config); err != nil {
		return nil, err
	}
	return result, nil
}
