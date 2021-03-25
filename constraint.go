package main

import (
	"github.com/go-git/go-billy/v5"
	"io/ioutil"
	"strings"
)

type Constraint interface {
	Evaluate(info billy.File) bool
}

type MultiConstraint struct {
	Constraints []Constraint
}

func (mc *MultiConstraint) Evaluate(info billy.File) bool {
	for _, c := range mc.Constraints {
		isGood := c.Evaluate(info)
		if !isGood {
			return false
		}
	}

	return true
}

type MustNotBeEmpty struct{}

func (_ *MustNotBeEmpty) Evaluate(file billy.File) bool {
	b, err := ioutil.ReadAll(file)

	if err != nil {
		Warning("✗ (%s) read error: %s", file.Name(), err)
		return false
	}
	if len(b) == 0 {
		Warning("✗ is empty: %s", file.Name())

		return false
	}

	return true
}

type Contains struct {
	Value string `yaml:"value"`
}

func (c *Contains) Evaluate(info billy.File) bool {
	b, err := ioutil.ReadAll(info)

	if err != nil {
		Warning("✗ (%s) read error: %s", info.Name(), err)
		return false
	}

	if strings.Contains(string(b), c.Value) {
		return true
	}

	Warning("✗ (%s) does not contain text: %s", info.Name(), c.Value)
	return false
}

type MustNotContain struct {
	Value string `yaml:"value"`
}

func (mnc *MustNotContain) Evaluate(info billy.File) bool {
	b, err := ioutil.ReadFile(info.Name())

	if err != nil {
		Warning("✗ (%s) read error: %s", info.Name(), err)
		return false
	}

	if strings.Contains(string(b), mnc.Value) {
		Warning("✗ (%s) contains text: %s", info.Name(), mnc.Value)
		return false
	}

	return true
}
