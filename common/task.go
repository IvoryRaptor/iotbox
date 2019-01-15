package common

import (
	"github.com/robfig/cron"
)

type ITask interface {
	cron.Job
	Config(kernel IKernel, config map[string]interface{}) error
	Work(channel IChannel)
}
