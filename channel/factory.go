package channel

import (
	"github.com/IvoryRaptor/iotbox/channel/mockchannel"
	"github.com/IvoryRaptor/iotbox/common"
)

func CreateChannel(ch chan common.ITask, config map[string]interface{}) (common.IChannel, error) {
	channelType := config["type"].(string)
	switch channelType {
	case "mock":
		mock := mockchannel.MockChannel{}
		if err := mock.Config(ch, config); err != nil {
			return nil, err
		}
		return &mock, nil
	}
	return nil, nil
}
