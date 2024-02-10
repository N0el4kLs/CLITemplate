package template

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/projectdiscovery/gologger"
)

//go:embed *.tpl
var f embed.FS

var (
	BANNER_TEMPLATE   = "banner.tpl"
	OPTION_TEMPLATE   = "option.tpl"
	RUNNER_TEMPLATE   = "runner.tpl"
	MAIN_TEMPLATE     = "main.tpl"
	MAKEFILE_TEMPLATE = "makefile.tpl"
	README_TEMPLATE   = "readme.tpl"
	IGNORE_TEMPLATE   = "gitignore.tpl"

	BANNER_FILEOUTPUT   = "runner/banner.go"
	OPTION_FILEOUTPUT   = "runner/option.go"
	RUNNER_FILEOUTPUT   = "runner/runner.go"
	MAKEFILE_FILEOUTPUT = "Makefile"
	README_FILEOUTPUT   = "README.md"
	IGNORE_FILEOUTPUT   = ".gitignore"

	Templates = []string{
		BANNER_TEMPLATE,
		OPTION_TEMPLATE,
		RUNNER_TEMPLATE,
		MAIN_TEMPLATE,
		MAKEFILE_TEMPLATE,
		README_TEMPLATE,
		IGNORE_TEMPLATE,
	}
)

var TemplateInstance *Template

type Template struct {
	ProjectName string
	GoModName   string
}

func Render(tplName string, workDir string) {
	gologger.Info().Msgf("Fill %s...\n", tplName)
	bContent, err := f.ReadFile(tplName)
	if err != nil {
		gologger.Fatal().Msgf("Could not read file: %s\n", err)
	}
	tmp, err := template.New(tplName).Parse(string(bContent))
	if err != nil {
		gologger.Fatal().Msgf("Could not parse template: %s\n", err)
	}

	var outputPath string
	if rlPath := getOutputPath(tplName); rlPath != "" {
		outputPath = filepath.Join(workDir, rlPath)
	} else {
		gologger.Fatal().Msgf("Could not get output path for template: %s\n", tplName)
	}
	f1, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, 0777)
	defer f1.Close()
	if err != nil {
		gologger.Fatal().Msgf("Could not Open file: %s\n", err)
	}
	err = tmp.Execute(f1, TemplateInstance)
	if err != nil {
		gologger.Fatal().Msgf("Could not execute template: %s\n", err)
	}
}

func getOutputPath(fileName string) string {
	switch fileName {
	case BANNER_TEMPLATE:
		return BANNER_FILEOUTPUT
	case OPTION_TEMPLATE:
		return OPTION_FILEOUTPUT
	case RUNNER_TEMPLATE:
		return RUNNER_FILEOUTPUT
	case MAIN_TEMPLATE:
		return fmt.Sprintf("cmd/%s/main.go", TemplateInstance.ProjectName)
	case MAKEFILE_TEMPLATE:
		return MAKEFILE_FILEOUTPUT
	case README_TEMPLATE:
		return README_FILEOUTPUT
	case IGNORE_TEMPLATE:
		return IGNORE_FILEOUTPUT
	}
	return ""
}
