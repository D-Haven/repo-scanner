package main

import (
	"encoding/json"
	"io/ioutil"
)

type Finding struct {
	Repository string   `json:"repository"`
	Errors     []string `json:"errors,flow"`
}

type Report struct {
	Successful []string  `json:"successful"`
	Findings   []Finding `json:"findings"`
}

func (report *Report) Write(path string) error {
	b, err := json.MarshalIndent(report, "", "  ")

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0666)
}
