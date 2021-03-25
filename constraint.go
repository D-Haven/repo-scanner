package main

import (
	"io/ioutil"
	"os"
	"strings"
)

type Constraint interface {
	Evaluate(info os.FileInfo) bool
}

type MultiConstraint struct {
	Constraints []Constraint
}

func (mc *MultiConstraint) Evaluate(info os.FileInfo) bool {
	for _, c := range mc.Constraints {
		isGood := c.Evaluate(info)
		if !isGood {
			return false
		}
	}

	return true
}

type MustNotBeEmpty struct{}

func (_ *MustNotBeEmpty) Evaluate(info os.FileInfo) bool {
	if info.Size() == 0 {
		Warning("✗ is empty: %s", info.Name())

		return false
	}

	return true
}

type Contains struct {
	Value string `yaml:"value"`
}

func (c *Contains) Evaluate(info os.FileInfo) bool {
	b, err := ioutil.ReadFile(info.Name())

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

func (mnc *MustNotContain) Evaluate(info os.FileInfo) bool {
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
