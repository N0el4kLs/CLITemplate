package main

import (
	"time"

	"CLITemplate/runner"

	"github.com/projectdiscovery/gologger"
)

func main() {
	options, err := runner.ParseOptions()
	if err != nil {
		panic(err)
	}
	r, err := runner.NewRunner(options)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	if err = r.Run(); err != nil {
		panic(err)
	}
	gologger.Info().Msgf("Task done,cost: %v\n", time.Since(start))
}
