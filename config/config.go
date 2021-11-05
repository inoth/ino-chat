package config

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	once sync.Once
	conf *Config
)

type Config struct {
	ServerPort string `yaml:"ServerPort"`
	MongoDB    struct {
		Host     string `yaml:"Host"`
		DataBase string `yaml:"DataBase"`
	} `yaml:"MongoDB"`
	Redis struct {
		Host   string `yaml:"Host"`
		Passwd string `yaml:"Passwd"`
	} `yaml:"Redis"`
	Nsq string `yaml:"Nsq"`
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
