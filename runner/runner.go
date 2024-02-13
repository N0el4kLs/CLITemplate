package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/N0el4kLs/CLITemplate/pkg/template"
	"github.com/projectdiscovery/gologger"
)

var (
	BaseDir, _ = os.Getwd()
	WorkDir    string
)

type Runner struct {
	options *Options
}

func NewRunner(option *Options) (*Runner, error) {
	runner := &Runner{}
	runner.options = option

	template.TemplateInstance = &template.Template{
		ProjectName: option.ProjectName,
		GoModName:   option.GoModName,
	}

	return runner, nil
}

func (r *Runner) Run() error {
	WorkDir = filepath.Join(BaseDir, r.options.ProjectName)

	// 1. create directories struct
	createDirectoryStruct(r.options.ProjectName)

	// 2. create go mod
	createGoMod(r.options.GoModName)

	// 3. fill template
	execTemplate(WorkDir)

	// 4. install dependencies
	installDependencies()
	return nil
}

func createGoMod(modeName string) {
	gologger.Info().Msgf("Create go mod...\n")
	c1 := exec.Command("go", "mod", "init", modeName)
	c1.Dir = WorkDir
	c1.Stdout = os.Stdout
	c1.Stderr = os.Stderr
	if err := c1.Run(); err != nil {
		gologger.Error().Msgf("Create go mod error:%s \n", err)
	}
}

func createDirectoryStruct(pjName string) {
	gologger.Info().Msgf("Create directories struct...\n")
	dirStructure := map[string][]string{
		fmt.Sprintf("cmd/%s", pjName): {"main.go"},
		"pkg":                         {},
		"runner":                      {"banner.go", "option.go", "runner.go"},
	}

	for dir, files := range dirStructure {
		dirPath := filepath.Join(WorkDir, dir)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			gologger.Fatal().Msgf("Failed to create directory: %s", dirPath)
		}

		for _, file := range files {
			filePath := filepath.Join(dirPath, file)
			_, err = os.Create(filePath)
			if err != nil {
				gologger.Fatal().Msgf("Failed to create file: %s", filePath)
			}
		}
	}

	files := []string{
		"Makefile",
		"README.md",
		".gitignore",
		"docker-compose.yaml",
	}
	for _, file := range files {
		filePath := filepath.Join(WorkDir, file)
		_, err := os.Create(filePath)
		if err != nil {
			gologger.Fatal().Msgf("Failed to create file: %s \n", filePath)
		}
	}
	gologger.Info().Msgf("Create directories struct done!\n")
}

func execTemplate(workDir string) {
	gologger.Info().Msgf("Fill template...\n")
	for _, tplName := range template.Templates {
		template.Render(tplName, workDir)
	}
}

func installDependencies() {
	gologger.Info().Msgf("Install dependences...\n")
	dependences := []string{
		"github.com/projectdiscovery/gologger",
	}
	for _, dependence := range dependences {
		gologger.Info().Msgf("Install dependence: %s \n", dependence)
		c := exec.Command("go", "get", dependence)
		c.Dir = WorkDir
		c.Stdout = os.Stdout
		if err := c.Run(); err != nil {
			gologger.Error().Msgf("Install dependence %s error: %s \n", dependence, err)
		}
	}
	c := exec.Command("go", "mod", "tidy")
	c.Dir = WorkDir
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		gologger.Error().Msgf("Go mod tidy error: %s\n", err)
	}
	gologger.Info().Msgf("Install dependencies done!\n")

	// format the go code
	gologger.Info().Msgf("Format the go code...\n")
	if err := exec.Command("go", "fmt", WorkDir).Run(); err != nil {
		gologger.Error().Msgf("Format the go code error: %s\n", err)
	}
	gologger.Info().Msgf("Format the go code done!\n")
}
