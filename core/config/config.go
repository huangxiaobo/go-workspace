package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"app"`

	Crawler struct {
		Tasks []struct {
			Project string `yaml:"project"`
			Url     string `yaml:"url"`
			Parser  string `yaml:"parser"`
		} `yaml:"tasks"`
	} `yaml:"crawler"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DbName   string `yaml:"dbname"`
	} `yaml:"database"`
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
}
