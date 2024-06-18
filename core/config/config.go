package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		Host string `yaml:"Host"`
		Port int    `yaml:"Port"`
	} `yaml:"app"`

}

var Conf Config

func LoadConfig(file string) {
	var err error
	conf, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(conf, &Conf)
	if err != nil {
		panic(err)
	}

	s, err := json.MarshalIndent(Conf, "", "    ")
	fmt.Println(string(s))
	if err != nil {
		panic("parse config file error")
	}
}
