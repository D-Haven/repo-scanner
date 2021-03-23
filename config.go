package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	WorkBranch string `yaml:"work-branch"`
    Username string `yaml:"username"`
	ApiTokenFile string `yaml:api-token`
	RequiredFiles []string `yaml:"required-files,flow"`
	Repositories []string `yaml:",flow"`
}

func ReadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := Config{}

	yaml.Unmarshal(b, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}