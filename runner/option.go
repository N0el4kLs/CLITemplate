package runner

import (
	"flag"
	"fmt"
	"os"
)

type Options struct {
	ProjectName string
	GoModName   string
}

func ParseOptions() (*Options, error) {
	ShowBanner()
	options := &Options{}

	flag.Usage = usage
	flag.StringVar(&options.ProjectName, "project", "", "Project name")
	flag.StringVar(&options.GoModName, "gomod", "", "Go mod name")
	flag.Parse()
	//flag.PrintDefaults()

	if options.ProjectName == "" {
		flag.Usage()
		os.Exit(0)
	}
	if options.GoModName == "" {
		options.GoModName = options.ProjectName
	}
	return options, nil
}

func usage() {
	comment := `
This is a tools to create my cli tool project template.
Usage:
  -gomod string
        Go mod name
  -project string
        Project name
`
	fmt.Println(comment)
}
