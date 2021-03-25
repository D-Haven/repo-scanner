package main

import (
	"strings"
)

type Constraint interface {
	Evaluate(filename string, b []byte) bool
}

type MultiConstraint struct {
	Constraints []Constraint
}

func (mc *MultiConstraint) Evaluate(filename string, b []byte) bool {
	for _, c := range mc.Constraints {
		isGood := c.Evaluate(filename, b)
		if !isGood {
			return false
		}
	}

	return true
}

type MustNotBeEmpty struct{}

func (_ *MustNotBeEmpty) Evaluate(filename string, b []byte) bool {
	if len(b) == 0 {
		Warning("✗ is empty: %s", filename)

		return false
	}

	return true
}

type Contains struct {
	Value string `yaml:"value"`
}

func (c *Contains) Evaluate(filename string, b []byte) bool {
	if strings.Contains(string(b), c.Value) {
		return true
	}

	Warning("✗ (%s) does not contain text: %s", filename, c.Value)
	return false
}

type MustNotContain struct {
	Value string `yaml:"value"`
}

func (mnc *MustNotContain) Evaluate(filename string, b []byte) bool {
	if strings.Contains(string(b), mnc.Value) {
		Warning("✗ (%s) contains text: %s", filename, mnc.Value)
		return false
	}

	return true
}
