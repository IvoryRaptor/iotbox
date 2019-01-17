package kernel

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strings"
)

func initModule(k *Kernel) error {
	files, _ := ioutil.ReadDir("./config/module")
	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
		data, err := ioutil.ReadFile(fmt.Sprintf("./config/module/%s", f.Name()))
		if err != nil {
			return err
		}
		var config map[string]interface{}
		if err := yaml.Unmarshal(data, &config); err != nil {
			return err
		}
		if m, err := k.CreateModule(config); err != nil {
			return err
		} else {
			c := make(chan common.ITask, 10)
			m.Start(c, m)
			fmt.Printf("Add Module [%s]\n", name)
			k.channel[name] = c
		}
	}
	return nil
}

func initTask(k *Kernel) error {
	files, _ := ioutil.ReadDir("./config/task")
	for _, f := range files {
		data, err := ioutil.ReadFile(fmt.Sprintf("./config/task/%s", f.Name()))
		if err != nil {
			return err
		}
		var config map[interface{}]interface{}
		if err := yaml.Unmarshal(data, &config); err != nil {
			return err
		}
		var cronTask common.ITask
		if cronTask, err = k.CreateTask(config); err != nil {
			return err
		}
		fmt.Printf("Add Corn Task [%s] %s\n", config["type"].(string), config["cron"].(string))
		k.JoinTask(config["cron"].(string), cronTask)
	}
	return nil
}

type Config struct {
	Channel []map[string]interface{}
}

func Boot() (*Kernel, error) {
	result := &Kernel{}
	result.cron = cron.New()
	result.channel = map[string]chan common.ITask{}
	if err := initModule(result); err != nil {
		return nil, err
	}
	if err := initTask(result); err != nil {
		return nil, err
	}
	return result, nil
}
