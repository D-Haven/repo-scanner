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

type FileTest struct {
	Name           string `yaml:"name"`
	RequireContent bool   `yaml:"not-empty,omitempty"`
	Includes       string `yaml:"includes,omitempty"`
	Excludes       string `yaml:"excludes,omitempty"`
}

type Config struct {
	WorkBranch    string     `yaml:"work-branch"`
	Auth          *RepoAuth  `yaml:"auth,omitempty"`
	RejectedFiles []string   `yaml:"rejected-files,flow"`
	RequiredFiles []FileTest `yaml:"required-files,flow"`
	Repositories  []string   `yaml:"repositories,flow"`
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

	// Default the username to "repo-scanner" if auth is specified but username is not
	if c.Auth != nil && len(c.Auth.Username) == 0 {
		c.Auth.Username = "repo-scanner"
	}

	return &c, nil
}

func (test *FileTest) Constraint() *MultiConstraint {
	constraint := MultiConstraint{}

	if test.RequireContent {
		constraint.Constraints = append(constraint.Constraints, &MustNotBeEmpty{})
	}

	if len(test.Includes) > 0 {
		constraint.Constraints = append(constraint.Constraints, &Contains{Value: test.Includes})
	}

	if len(test.Excludes) > 0 {
		constraint.Constraints = append(constraint.Constraints, &MustNotContain{Value: test.Excludes})
	}

	return &constraint
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
