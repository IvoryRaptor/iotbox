package corn

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

type Corn struct {
	common.AModule
	cron *cron.Cron
}

func (m *Corn) GetName() string {
	return "core"
}

func (m *Corn) Config(kernel common.IKernel, config map[string]interface{}) error {
	m.cron = cron.New()
	tasksPath := config["tasks"].(string)
	files, _ := ioutil.ReadDir(tasksPath)
	for _, f := range files {
		if path.Ext(f.Name()) == ".yaml" {
			data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", tasksPath, f.Name()))
			if err != nil {
				return err
			}
			var config map[string]interface{}
			if err := yaml.Unmarshal(data, &config); err != nil {
				return err
			}
			var cronTask common.ITask
			if cronTask, err = kernel.CreateTask(config); err != nil {
				return err
			}
			fmt.Printf("Add Corn Task [%s] %s\n", config["type"].(string), config["cron"].(string))
			m.cron.AddJob(config["cron"].(string), cronTask)
		}
	}
	m.cron.Start()
	return nil
}

func (m *Corn) Send(_ common.ITask, packet common.Packet) chan common.Packet {
	return m.Response
}

func CreateCore() *Corn {
	return &Corn{}
}
