package common

import (
	"errors"
	"github.com/robfig/cron"
)

type WorkState int

const (
	Running WorkState = iota
	Complete
	Failed
)

type ITask interface {
	cron.Job
	Config(kernel IKernel, config map[interface{}]interface{}) error
	Work(channel IModule) (WorkState, error)
}

type IHandlerTask interface {
	ITask
	Clone() IHandlerTask
	SetPacket(packet Packet) IHandlerTask
}

type IOwnerTask interface {
	ITask
	GetRequest() Packet
	SetOwner(owner IModule)
}

type ATask struct {
	target      string
	kernel      IKernel
	CurrentWork func(channel IModule) (WorkState, error)
	OtherConfig func(kernel IKernel, config map[interface{}]interface{}) error
}

func (t *ATask) Work(channel IModule) (WorkState, error) {
	if t.CurrentWork != nil {
		return t.CurrentWork(channel)
	}
	return Failed, errors.New("CurrentWork is nil")
}

func (t *ATask) Run() {
	t.kernel.GetModule(t.target) <- t
}

func (t *ATask) Config(kernel IKernel, config map[interface{}]interface{}) error {
	t.target = config["target"].(string)
	t.kernel = kernel
	if t.OtherConfig != nil {
		return t.OtherConfig(kernel, config)
	}
	return nil
}

func (t *ATask) SetCurrentWork(work func(channel IModule) (WorkState, error)) *ATask {
	t.CurrentWork = work
	return t
}

func (t *ATask) SetOtherConfig(config func(kernel IKernel, config map[interface{}]interface{}) error) *ATask {
	t.OtherConfig = config
	return t
}
