package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Wiki struct {
		Name string `yaml:"name"`
		Desc string `yaml:"desc"`
		Icon string `yaml:"icon"`
	} `yaml:"Wiki"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"Server"`
	Mongo struct {
		Conn string `yaml:"conn"`
		Db   string `yaml:"db"`
	} `yaml:"Mongo"`
}

var config Config

func LoadConfig() (*Config, error) {
	config_file, err := ioutil.ReadFile("./configs/config.yaml")
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(config_file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func GetConfig() *Config {
	return &config
}
