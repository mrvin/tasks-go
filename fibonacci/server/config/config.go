package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HostGRPC string `yaml:"hostGRPC"`
	PortGRPC int    `yaml:"portGRPC"`
	HostHTTP string `yaml:"hostHTTP"`
	PortHTTP int    `yaml:"portHTTP"`
	DB       struct {
		Host     string `yaml:"hostDB"`
		Port     int    `yaml:"portDB"`
		Password string `yaml:"passwordDB"`
		NameDB   int    `yaml:"nameDB"`
	} `yaml:"db"`
}

func (conf *Config) Parse(configPath string) error {
	configYml, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading %s error: %w", configPath, err)
	}

	err = yaml.Unmarshal(configYml, conf)
	if err != nil {
		return fmt.Errorf("can't unmarshal %s: %w", configPath, err)
	}

	return nil
}
