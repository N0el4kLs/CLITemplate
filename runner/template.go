package runner

// will be replaced by the template in next version
var (
	MakefileTemplate string = `# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test
GOFLAGS := -v
LDFLAGS := -s -w

ifneq ($(shell go env GOOS),darwin)
LDFLAGS := -extldflags "-static"
endif

all: build
build:
	$(GOBUILD) $(GOFLAGS) -ldflags '$(LDFLAGS)' -o "%s" cmd/%s/main.go
tidy:
	$(GOMOD) tidy
`
	BannerTemplate string = `package runner

import "fmt"

const version = "v0.0.1"

func ShowBanner() {
	//http://www.network-science.de/ascii/  smslant
	var banner = ''
	fmt.Printf(banner, version)
}
`
	OptionTemplate string = `package runner

type Options struct {
}

func ParseOptions() (*Options,error) {
	options := &Options{}
	return options,nil
}
`
	RunnerTemplate string = `package runner

func NewRunner(option *Options) (*Runner, error) {
	runner := &Runner{}
	runner.options = option

	return runner, nil
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Run() {

}

func (r *Runner) Close() {

}
`
	MainTemplate string = `package main

import (
	"time"

	"%s/runner"

	"github.com/projectdiscovery/gologger"
)

func main() {
	options, err := runner.ParseOptions()
	if err != nil {
		gologger.Fatal().Msgf("Could not parse options: %%s\n", err)
	}
	r, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %%s\n", err)
	}

	start := time.Now()
	if err = r.Run(); err != nil {
		panic(err)
	}
	gologger.Info().Msgf("Task done,cost: %%v\n", time.Since(start))
}
`
	ReadmeTemplate = `<h2 align="center">%s</h2>

## 介绍

`
	IgnoreTemplate = `.idea/
.vscode
*_test.go
`
)
