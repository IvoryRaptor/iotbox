package report

import (
	"github.com/IvoryRaptor/iotbox/common"
	"time"
)

type Report struct {
	common.ATask
	index  int
	packet common.Packet
}

func (d *Report) ReportConfig(kernel common.IKernel, config map[interface{}]interface{}) error {
	d.SetCurrentWork(d.StartWork)
	return nil
}

func (d *Report) Clone() common.ICloneTask {
	result := &Report{
		ATask: d.ATask,
	}
	result.SetCurrentWork(result.StartWork)
	return result
}

func (d *Report) SetPacket(packet common.Packet) (common.ICloneTask, error) {
	d.packet = packet
	return d, nil
}

func (d *Report) StartWork(module common.IModule) (common.WorkState, error) {
	ch := module.Send(d, d.packet)
	module.Read(ch, time.Second*3)
	return common.Complete, nil
}

func Create() *Report {
	result := &Report{}
	result.SetOtherConfig(result.ReportConfig)
	return result
}
