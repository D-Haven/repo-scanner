package main

import (
	"fmt"
	"strings"
)

type Constraint interface {
	Evaluate(filename string, b []byte) (bool, string)
}

type MultiConstraint struct {
	Constraints []Constraint
}

func (mc *MultiConstraint) Evaluate(filename string, b []byte) (bool, string) {
	for _, c := range mc.Constraints {
		isGood, message := c.Evaluate(filename, b)
		if !isGood {
			return false, message
		}
	}

	return true, ""
}

type MustNotBeEmpty struct{}

func (_ *MustNotBeEmpty) Evaluate(filename string, b []byte) (bool, string) {
	if len(b) == 0 {
		return false, fmt.Sprintf("empty: %s", filename)
	}

	return true, ""
}

type Contains struct {
	Value string `yaml:"value"`
}

func (c *Contains) Evaluate(filename string, b []byte) (bool, string) {
	if strings.Contains(string(b), c.Value) {
		return true, ""
	}

	return false, fmt.Sprintf("(%s) does not contain text: %s", filename, c.Value)
}

type MustNotContain struct {
	Value string `yaml:"value"`
}

func (mnc *MustNotContain) Evaluate(filename string, b []byte) (bool, string) {
	if strings.Contains(string(b), mnc.Value) {
		return false, fmt.Sprintf("(%s) contains text: %s", filename, mnc.Value)
	}

	return true, ""
}
