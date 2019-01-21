package common

import (
	"fmt"
	"time"
)

type IModule interface {
	GetName() string
	Config(kernel IKernel, config map[string]interface{}) error
	Send(handle ITask, packet Packet) chan Packet
	Start(this IModule)
	GetTaskQueue() chan ITask
	Read(ch chan Packet, timeOut time.Duration) Packet
}

type AModule struct {
	Response   chan Packet
	taskQueue  chan ITask
	QueryCount int
}

func (m *AModule) Read(ch chan Packet, timeOut time.Duration) Packet {
	select {
	case res := <-ch:
		return res
	case <-time.After(timeOut):
		return nil
	}
}

func (m *AModule) GetTaskQueue() chan ITask {
	return m.taskQueue
}

func (m *AModule) Start(this IModule) {
	if m.QueryCount <= 0 {
		m.QueryCount = 10
	}
	m.Response = make(chan Packet)
	m.taskQueue = make(chan ITask, m.QueryCount)
	go func() {
		for {
			task := <-m.GetTaskQueue()
			if state, err := task.Work(this); err != nil {
				fmt.Printf("Task Work Error %s\n", err.Error())
			} else {
				switch state {
				case Complete:

				case Failed:

				case Running: //未完成任务，重新回归队列
					m.taskQueue <- task
				}
			}
		}
	}()
}
