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

func (mc *MultiConstraint) Evaluate(file billy.File) bool {
	for _, c := range mc.Constraints {
		_, err := file.Seek(0, 0)
		if err != nil {
			Warning("✗ (%s) seek error: %s", file.Name(), err)
			return false
		}

		isGood := c.Evaluate(file)
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

func (c *Contains) Evaluate(file billy.File) bool {
	b, err := ioutil.ReadAll(file)

	if err != nil {
		Warning("✗ (%s) read error: %s", file.Name(), err)
		return false
	}

	if strings.Contains(string(b), c.Value) {
		return true
	}

	Warning("✗ (%s) does not contain text: %s", file.Name(), c.Value)
	return false
}

type MustNotContain struct {
	Value string `yaml:"value"`
}

func (mnc *MustNotContain) Evaluate(file billy.File) bool {
	b, err := ioutil.ReadFile(file.Name())

	if err != nil {
		Warning("✗ (%s) read error: %s", file.Name(), err)
		return false
	}

	if strings.Contains(string(b), mnc.Value) {
		Warning("✗ (%s) contains text: %s", file.Name(), mnc.Value)
		return false
	}

	return true
}
