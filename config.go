package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"strings"
)

type RepoAuth struct {
	Mode         string `yaml:"mode"`
	Username     string `yaml:"username,omitempty"`
	ApiTokenFile string `yaml:"api-token"`
}

type Config struct {
	WorkBranch    string    `yaml:"work-branch"`
	Auth          *RepoAuth `yaml:"auth,omitempty"`
	RequiredFiles []string  `yaml:"required-files,flow"`
	Repositories  []string  `yaml:"repositories,flow"`
}

func ReadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := Config{}

	err = yaml.Unmarshal(b, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (auth *RepoAuth) Transform(content string) (string, error) {
	if strings.EqualFold(auth.Mode, "url") {
		if len(content) < 8 {
			return content, nil
		}

		scheme := content[:8]

		if !strings.EqualFold(scheme, "https://") {
			return "", errors.New("URL authentication only supported for HTTPS")
		}

		var path = content[8:]

		if len(auth.Username) == 0 {
			auth.Username = "repo-scanner"
		}

		b, err := ioutil.ReadFile(auth.ApiTokenFile)
		CheckIfError(err)
		token := strings.TrimSpace(string(b))

		path = fmt.Sprintf("%s:%s@%s", auth.Username, url.PathEscape(token), path)

		return scheme + path, nil
	}

	return content, nil
}
