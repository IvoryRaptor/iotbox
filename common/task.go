package common

import (
	"errors"
	"fmt"
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
	Config(kernel IKernel, config map[string]interface{}) error
	Work(channel IModule) (WorkState, error)
}

type ICloneTask interface {
	ITask
	Clone() ICloneTask
	SetPacket(packet Packet) (ICloneTask, error)
}

type IOwnerTask interface {
	ITask
	GetRequest() Packet
	SetOwner(owner IModule)
}

type ATask struct {
	target        string
	kernel        IKernel
	CurrentModule IModule
	CurrentWork   func(channel IModule) (WorkState, error)
	OtherConfig   func(kernel IKernel, config map[string]interface{}) error
}

func (t *ATask) Work(module IModule) (WorkState, error) {
	if t.CurrentModule != nil && t.CurrentModule == module {
		return Failed, errors.New(fmt.Sprintf("module [%s] will call task,but module [%s] is calling", module.GetName(), t.CurrentModule.GetName()))
	}
	if t.CurrentWork == nil {
		return Failed, errors.New("CurrentWork is nil")
	}
	state, err := t.CurrentWork(module)
	switch state {
	case Complete:
		t.CurrentModule = nil
	case Failed:
		t.CurrentModule = nil
	}
	return state, err
}

func (t *ATask) Run() {
	t.JoinQueue(t.target)
}

func (t *ATask) JoinQueue(module string) *ATask {
	t.kernel.GetModule(module) <- t
	return t
}

func (t *ATask) Config(kernel IKernel, config map[string]interface{}) error {
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

func (t *ATask) SetOtherConfig(config func(kernel IKernel, config map[string]interface{}) error) *ATask {
	t.OtherConfig = config
	return t
}
