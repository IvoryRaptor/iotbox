package rootcloud

import (
	"github.com/IvoryRaptor/iotbox/common"
	"time"
	"encoding/json"
	"log"
)

// RootCloud 根云上报
type RootCloud struct
{
	common.ATask
	packet common.Packet
}

func (rc *RootCloud) otherConfig(kernel common.IKernel, config map[string]interface{}) error {
	rc.SetCurrentWork(rc.StartWork)
	return nil
}

// Clone 复制task
func (rc *RootCloud) Clone() common.ICloneTask {
	result := &RootCloud{
		ATask: rc.ATask,
	}
	result.SetCurrentWork(result.StartWork)
	return result
}

// SetPacket 设置上报包
func (rc *RootCloud) SetPacket(packet common.Packet) (common.ICloneTask, error) {
	rc.packet = packet
	packet["test"] = make([]byte, 10)
	return rc, nil
}

// StartWork 进行上报
func (rc *RootCloud) StartWork(module common.IModule) (common.WorkState, error) {
	reportMap := make(map[string]interface{})
	reportMap[rc.packet["Name"].(string)] = rc.packet["ConversionValue"]
	reportData, err := json.Marshal(reportMap)
	if err != nil {
		log.Fatalf("[rootcloud]===> StartWork %s", err)
		return common.Complete, nil
	}
	ch := module.Send(rc, common.Packet{"value": reportData})
	module.Read(ch, time.Second*3)
	return common.Complete, nil
}

// Create 创建rootcloud task
func Create() *RootCloud {
	result := &RootCloud{}
	result.SetOtherConfig(result.otherConfig)
	return result
}