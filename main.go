package main

import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
	log.Printf("INFO: %s", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
	log.Printf("WARN: %s", fmt.Sprintf(format, args...))
}

func LogError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	log.Printf("ERR:  %s", err)
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	LogError(err)
	os.Exit(1)
}

func main() {
	var configFile = "scan-repos.yaml"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	c, err := ReadConfig(configFile)
	CheckIfError(err)

	ext := path.Ext(configFile)
	logfile := configFile[0:len(configFile)-len(ext)] + ".report"
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	CheckIfError(err)
	log.SetOutput(file)

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

	for _, repo := range c.Repositories {
		// Clones the given repository in memory, creating the remote, the local
		// branches and fetching the objects, exactly as:
		Info("git clone --single-branch %s %s", c.WorkBranch, repo)

		var transformedUrl = repo

		if c.Auth != nil {
			transformedUrl, err = c.Auth.Transform(repo)
			if err != nil {
				LogError(err)
				continue
			}
		}

		r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
			URL:           transformedUrl,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", c.WorkBranch)),
			SingleBranch:  true,
			Auth:          auth,
		})
		LogError(err)
		if err != nil {
			continue
		}

		wt, err := r.Worktree()
		CheckIfError(err)

		var isGood = true

		for _, file := range c.RequiredFiles {
			info, err := wt.Filesystem.Open(file.Name)

			if err != nil {
				isGood = false
				Warning("✗ missing: %s", file.Name)
				continue
			}

			isGood = isGood && file.Constraint().Evaluate(info)
			err = info.Close()
			if err != nil {
				Warning("✗ can't close %s: %s", info.Name(), err)
			}
		}

		if isGood {
			Info("✓ has all required files")
		}
	}
}
