package kernel

import (
	"fmt"
	"github.com/IvoryRaptor/iotbox/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strings"
)

func initModule(k *Kernel) error {
	files, _ := ioutil.ReadDir("./config/module")
	for _, f := range files {
		if path.Ext(f.Name()) == ".yaml" {
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
	}
	return nil
}

type Config struct {
	Channel []map[string]interface{}
}

func Boot() (*Kernel, error) {
	result := &Kernel{}
	result.channel = map[string]chan common.ITask{}
	if err := initModule(result); err != nil {
		return nil, err
	}
	return result, nil
}
