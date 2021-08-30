package config

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	// 房间列表
	// ROOMS = "INOCHAT:ROOMS"
	// 房间详情
	ROOMINFO = "INOCHAT:ROOMINFO:"
	// 房间内成员
	ROOMMEMBERS = "INOCHAT:ROOM:"
)

var (
	once sync.Once
	conf *Config
)

type Config struct {
	ServerPort string `yaml:"ServerPort"`
	Redis      struct {
		Host   string `yaml:"Host"`
		Passwd string `yaml:"Passwd"`
	} `yaml:"Redis"`
	MongoDB string `yaml:"MongoDB"`
	Consul  string `yaml:"Consul"`
	Nsq     string `yaml:"Nsq"`
}

func Instance() *Config {
	once.Do(func() {
		conf = &Config{}
		yamlFile, err := ioutil.ReadFile(selectConfigPath(nil))
		if err != nil {
			logrus.Errorf("%v", err)
			os.Exit(1)
		}
		err = yaml.Unmarshal(yamlFile, conf)
		if err != nil {
			logrus.Errorf("%v", err)
			os.Exit(1)
		}
	})
	return conf
}

// func init() {
// 	yamlFile, err := ioutil.ReadFile(selectConfigPath(nil))
// 	if err != nil {
// 		log.Errorf("%v", err)
// 		os.Exit(1)
// 	}
// 	err = yaml.Unmarshal(yamlFile, Instance())
// 	if err != nil {
// 		log.Errorf("%v", err)
// 		os.Exit(1)
// 	}
// }

func selectConfigPath(path []string) string {
	if len(path) > 0 {
		return path[0]
	} else {
		e := os.Getenv("GORUNEVN")
		if len(e) > 0 {
			return "config." + e + ".yaml"
		}
		return "config.yaml"
	}
}
