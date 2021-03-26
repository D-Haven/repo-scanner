package main

import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m\t%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func main() {
	var configFile = "scan-repos.yaml"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	c, err := ReadConfig(configFile)
	CheckIfError(err)

	// Support Github and Bitbucket API Tokens
	var auth *http.BasicAuth = nil

	if c.Auth != nil && strings.EqualFold(c.Auth.Mode, "basic") {
		auth := &http.BasicAuth{
			Username: "repo-scanner",
		}

		if len(c.Auth.Username) > 0 {
			auth.Username = c.Auth.Username
		}

		b, err := ioutil.ReadFile(c.Auth.ApiTokenFile)
		CheckIfError(err)
		auth.Password = strings.TrimSpace(string(b))
	}

	var wg = sync.WaitGroup{}
	wg.Add(len(c.Repositories))
	report := Report{}

	for _, repo := range c.Repositories {
		go func(repo string) {
			defer wg.Done()
			// Clones the given repository in memory, creating the remote, the local
			// branches and fetching the objects, exactly as:
			Info("git clone --single-branch %s %s", c.WorkBranch, repo)
			finding := Finding{
				Repository: repo,
			}

			var transformedUrl = repo

			if c.Auth != nil {
				transformedUrl, err = c.Auth.Transform(repo)
				if err != nil {
					finding.Errors = append(finding.Errors, err.Error())
					report.Findings = append(report.Findings, finding)
					return
				}
			}

			r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
				URL:           transformedUrl,
				ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", c.WorkBranch)),
				SingleBranch:  true,
				Auth:          auth,
			})
			if err != nil {
				finding.Errors = append(finding.Errors, err.Error())
				report.Findings = append(report.Findings, finding)
				return
			}

			wt, err := r.Worktree()
			if err != nil {
				finding.Errors = append(finding.Errors, err.Error())
				report.Findings = append(report.Findings, finding)
				return
			}

			var isGood = true
			for _, rejected := range c.RejectedFiles {
				_, err := wt.Filesystem.Stat(rejected)

				if err != nil {
					finding.Errors = append(finding.Errors, fmt.Sprintf("has rejected file: %s", rejected))
				}
			}

			for _, expected := range c.RequiredFiles {
				file, err := wt.Filesystem.Open(expected.Name)

				if err != nil {
					isGood = false
					message := fmt.Sprintf("missing: %s", expected.Name)
					finding.Errors = append(finding.Errors, message)
					continue
				}

				b, err := ioutil.ReadAll(file)
				if err != nil {
					isGood = false
					message := fmt.Sprintf("(%s) read error: %s", file.Name(), err)
					finding.Errors = append(finding.Errors, message)
					continue
				}

				err = file.Close()
				if err != nil {
					isGood = false
					message := fmt.Sprintf("can't close %s: %s", file.Name(), err)
					finding.Errors = append(finding.Errors, message)
					continue
				}

				passed, message := expected.Constraint().Evaluate(file.Name(), b)

				if !passed {
					finding.Errors = append(finding.Errors, message)
				}

				isGood = isGood && passed
			}

			if isGood {
				report.Successful = append(report.Successful, repo)
			} else {
				report.Findings = append(report.Findings, finding)
			}
		}(repo)
	}

	wg.Wait()

	Info("... Discovered %d findings, read the report for details", len(report.Findings))

	ext := path.Ext(configFile)
	reportFile := configFile[0:len(configFile)-len(ext)] + ".report"
	err = report.Write(reportFile)
	CheckIfError(err)
}
