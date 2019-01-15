package kernel

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/channel"
	"github.com/IvoryRaptor/iotbox/common"
	"github.com/IvoryRaptor/iotbox/task"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strings"
)

func initConfig(k *Kernel) error {
	data, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		return err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	for _, chanConfig := range config.Channel {
		c := make(chan common.ITask, 10)
		if _, err := channel.CreateChannel(c, chanConfig); err != nil {
			return err
		}
		k.channel[chanConfig["name"].(string)] = c
	}
	return nil
}

func initChannel(k *Kernel) error {
	files, _ := ioutil.ReadDir("./config/channel")
	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
		data, err := ioutil.ReadFile(fmt.Sprintf("./config/channel/%s", f.Name()))
		if err != nil {
			return err
		}
		var config map[string]interface{}
		if err := yaml.Unmarshal(data, &config); err != nil {
			return err
		}
		c := make(chan common.ITask, 10)
		if _, err := channel.CreateChannel(c, config); err != nil {
			return err
		}
		fmt.Printf("Add Channel [%s]\n", name)
		k.channel[name] = c
	}
	return nil
}

func initTask(kernel *Kernel) error {
	files, _ := ioutil.ReadDir("./config/task")
	for _, f := range files {
		data, err := ioutil.ReadFile(fmt.Sprintf("./config/task/%s", f.Name()))
		if err != nil {
			return err
		}
		var config map[string]interface{}
		if err := yaml.Unmarshal(data, &config); err != nil {
			return err
		}
		if err := task.CreateTask(kernel, config); err != nil {
			return err
		}
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
	if err := initConfig(result); err != nil {
		return nil, err
	}
	if err := initChannel(result); err != nil {
		return nil, err
	}
	if err := initTask(result); err != nil {
		return nil, err
	}
	return result, nil
}
