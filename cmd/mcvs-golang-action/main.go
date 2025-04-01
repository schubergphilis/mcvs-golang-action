package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Vars struct {
	RemoteURL     string `yaml:"REMOTE_URL"`
	RemoteURLRef  string `yaml:"REMOTE_URL_REF"`
	RemoteURLRepo string `yaml:"REMOTE_URL_REPO"`
}

type Includes struct {
	Remote struct {
		Taskfile string `yaml:"taskfile"`
		Vars     struct {
			GolangciLintRunTimeoutMinutes int `yaml:"GOLANGCI_LINT_RUN_TIMEOUT_MINUTES"`
		} `yaml:"vars"`
	} `yaml:"remote"`
}

type Config struct {
	Version  int      `yaml:"version"`
	Vars     Vars     `yaml:"vars"`
	Includes Includes `yaml:"includes"`
}

type Workflow struct {
	Name        string            `yaml:"name"`
	On          string            `yaml:"on"`
	Permissions map[string]string `yaml:"permissions"`
	Jobs        map[string]Job    `yaml:"jobs"`
}

type Job struct {
	Strategy Strategy `yaml:"strategy"`
	RunsOn   string   `yaml:"runs-on"`
	Env      Env      `yaml:"env"`
	Steps    []Step   `yaml:"steps"`
}

type Strategy struct {
	Matrix Matrix `yaml:"matrix"`
}

type Matrix struct {
	Args []map[string]interface{} `yaml:"args"`
}

type Env struct {
	TaskXRemoteTaskfiles int    `yaml:"TASK_X_REMOTE_TASKFILES"`
	TestTimeout          string `yaml:"test-timeout"`
}

type Step struct {
	Uses string                 `yaml:"uses,omitempty"`
	With map[string]interface{} `yaml:"with,omitempty"`
}

type Dependabot struct {
	Version int      `yaml:"version"`
	Updates []Update `yaml:"updates"`
}

type Update struct {
	PackageEcosystem string           `yaml:"package-ecosystem"`
	Directory        string           `yaml:"directory"`
	Schedule         Schedule         `yaml:"schedule"`
	Groups           map[string]Group `yaml:"groups"`
}

type Schedule struct {
	Interval string `yaml:"interval"`
}

type Group struct {
	Patterns []string `yaml:"patterns"`
}

func createGolangWorkflow(filePath string) error {
	golangUnitTestsExclusions := `
\(cmd\/some-app\|internal\/app\/some-app\)`

	workflow := Workflow{
		Name: "Golang",
		On:   "push",
		Permissions: map[string]string{
			"contents": "read",
			"packages": "read",
		},
		Jobs: map[string]Job{
			"mcvs-golang-action": {
				Strategy: Strategy{
					Matrix: Matrix{
						Args: []map[string]interface{}{
							{
								"release-application-name": "mcvs-golang-action-bootstrap",
								"release-architecture":     "arm64",
								"release-dir":              "./cmd/mcvs-golang-action",
								"release-type":             "binary",
							},
							{"testing-type": "lint", "build-tags": ""},
							{"testing-type": "security-golang-modules"},
							{"testing-type": "security-grype"},
							{"testing-type": "security-trivy", "security-trivyignore": ""},
							{"testing-type": "unit"},
						},
					},
				},
				RunsOn: "ubuntu-22.04",
				Env: Env{
					TaskXRemoteTaskfiles: 1,
					TestTimeout:          "10m0s",
				},
				Steps: []Step{
					{Uses: "actions/checkout@v4.1.1"},
					{
						Uses: "schubergphilis/mcvs-golang-action@v1.0.1",
						With: map[string]any{
							"build-tags":                   "${{ matrix.args.build-tags }}",
							"code-coverage-expected":       0.0,
							"code-coverage-timeout":        "${{ env.test-timeout }}",
							"golang-unit-tests-exclusions": golangUnitTestsExclusions,
							"release-architecture":         "${{ matrix.args.release-architecture }}",
							"release-dir":                  "${{ matrix.args.release-dir }}",
							"release-type":                 "${{ matrix.args.release-type }}",
							"security-trivyignore":         "${{ matrix.args.security-trivyignore }}",
							"testing-type":                 "${{ matrix.args.testing-type }}",
							"test-timeout":                 "${{ env.test-timeout }}",
							"token":                        "${{ secrets.GITHUB_TOKEN }}",
						},
					},
				},
			},
		},
	}

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	err := encoder.Encode(&workflow)
	if err != nil {
		fmt.Printf("Error encoding YAML: %v\n", err)
	}
	encoder.Close()
	yamlOutput := buffer.String()
	yamlOutput = strings.ReplaceAll(yamlOutput, "'", "\"")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString("---\n" + yamlOutput)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}

	return nil
}

func createDependabotConfig(filePath string) error {
	dependabot := Dependabot{
		Version: 2,
		Updates: []Update{
			{
				PackageEcosystem: "docker",
				Directory:        "/",
				Schedule: Schedule{
					Interval: "weekly",
				},
				Groups: map[string]Group{
					"docker-all": {
						Patterns: []string{"*"},
					},
				},
			},
			{
				PackageEcosystem: "github-actions",
				Directory:        "/",
				Schedule: Schedule{
					Interval: "weekly",
				},
				Groups: map[string]Group{
					"github-actions-all": {
						Patterns: []string{"*"},
					},
				},
			},
			{
				PackageEcosystem: "gomod",
				Directory:        "/",
				Schedule: Schedule{
					Interval: "weekly",
				},
				Groups: map[string]Group{
					"gomod-all": {
						Patterns: []string{"*"},
					},
				},
			},
		},
	}

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	err := encoder.Encode(&dependabot)
	if err != nil {
		fmt.Printf("Error encoding YAML: %v\n", err)
	}
	encoder.Close()
	yamlOutput := buffer.String()
	yamlOutput = strings.ReplaceAll(yamlOutput, "'", "\"")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString("---\n" + yamlOutput)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}

	return nil
}

func bla() {
	config := Config{
		Version: 3,
		Vars: Vars{
			RemoteURL:     "https://raw.githubusercontent.com",
			RemoteURLRef:  "v1.0.1",
			RemoteURLRepo: "schubergphilis/mcvs-golang-action",
		},
		Includes: Includes{
			Remote: struct {
				Taskfile string `yaml:"taskfile"`
				Vars     struct {
					GolangciLintRunTimeoutMinutes int `yaml:"GOLANGCI_LINT_RUN_TIMEOUT_MINUTES"`
				} `yaml:"vars"`
			}{
				Taskfile: `{{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/Taskfile.yml`,
				Vars: struct {
					GolangciLintRunTimeoutMinutes int `yaml:"GOLANGCI_LINT_RUN_TIMEOUT_MINUTES"`
				}{
					GolangciLintRunTimeoutMinutes: 5,
				},
			},
		},
	}

	filename := "Taskfile.yml"

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	err := encoder.Encode(&config)
	if err != nil {
		fmt.Printf("Error encoding YAML: %v\n", err)
	}
	encoder.Close()
	yamlOutput := buffer.String()
	yamlOutput = strings.ReplaceAll(yamlOutput, "'", "\"")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
	}
	defer file.Close()

	modifiedYamlString := strings.Replace(yamlOutput, `"{{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/Taskfile.yml"`, `>-
      {{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/Taskfile.yml`, 1)

	_, err = file.WriteString("---\n" + modifiedYamlString)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func gitignore() {
	filename := ".gitignore"
	ignoreEntry := ".task"

	// Check if the .gitignore file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// If the file does not exist, create it
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Failed to create %s: %v\n", filename, err)
			return
		}
		defer file.Close()
		fmt.Printf("%s created successfully.\n", filename)
	}

	// Read the contents of .gitignore
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read %s: %v\n", filename, err)
		return
	}

	// Check if the ignoreEntry is already present
	lines := strings.Split(string(data), "\n")
	found := false
	for _, line := range lines {
		if line == ignoreEntry {
			found = true
			break
		}
	}

	// If not found, append the ignoreEntry to the file
	if !found {
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Printf("Failed to open %s: %v\n", filename, err)
			return
		}
		defer file.Close()

		if _, err := file.WriteString(ignoreEntry + "\n"); err != nil {
			fmt.Printf("Failed to write to %s: %v\n", filename, err)
			return
		}

		fmt.Printf("Added '%s' to %s.\n", ignoreEntry, filename)
	} else {
		fmt.Printf("'%s' is already present in %s.\n", ignoreEntry, filename)
	}
}

type DepGuard struct {
	Rules struct {
		Main struct {
			Files []string `yaml:"files"`
			Allow []string `yaml:"allow"`
			Deny  []struct {
				Pkg  string `yaml:"pkg"`
				Desc string `yaml:"desc"`
			} `yaml:"deny"`
		} `yaml:"main"`
	} `yaml:"rules"`
}

type Exclusions struct {
	Generated string   `yaml:"generated"`
	Presets   []string `yaml:"presets"`
	Rules     []Rule   `yaml:"rules"`
	Paths     []string `yaml:"paths"`
}

type Settings struct {
	DepGuard DepGuard `yaml:"depguard"`
}

type FormatterSettings struct {
	Enable     []string `yaml:"enable"`
	Exclusions struct {
		Generated string   `yaml:"generated"`
		Paths     []string `yaml:"paths"`
	} `yaml:"exclusions"`
}

type Rule struct {
	Linters []string `yaml:"linters"`
	Path    string   `yaml:"path"`
	Text    string   `yaml:"text"`
}

type Config2 struct {
	Version string `yaml:"version"`
	Linters struct {
		Default    string     `yaml:"default"`
		Disable    []string   `yaml:"disable"`
		Settings   Settings   `yaml:"settings"`
		Exclusions Exclusions `yaml:"exclusions"`
	} `yaml:"linters"`
	Formatters FormatterSettings `yaml:"formatters"`
}

func golangCILint() {
	config := Config2{
		Version: "2",
		Linters: struct {
			Default    string     `yaml:"default"`
			Disable    []string   `yaml:"disable"`
			Settings   Settings   `yaml:"settings"`
			Exclusions Exclusions `yaml:"exclusions"`
		}{
			Default: "all",
			Disable: []string{
				"exhaustruct", "forcetypeassert", "gochecknoglobals",
				"gocritic", "lll", "nestif", "noctx", "nonamedreturns", "perfsprint",
				"testifylint", "testpackage", "varnamelen", "wrapcheck",
			},
			Settings: Settings{
				DepGuard: DepGuard{
					Rules: struct {
						Main struct {
							Files []string `yaml:"files"`
							Allow []string `yaml:"allow"`
							Deny  []struct {
								Pkg  string `yaml:"pkg"`
								Desc string `yaml:"desc"`
							} `yaml:"deny"`
						} `yaml:"main"`
					}{
						Main: struct {
							Files []string `yaml:"files"`
							Allow []string `yaml:"allow"`
							Deny  []struct {
								Pkg  string `yaml:"pkg"`
								Desc string `yaml:"desc"`
							} `yaml:"deny"`
						}{
							Files: []string{"!**/*_a _file.go"},
							Allow: []string{
								"$gostd", "github.com/aws/aws-lambda-go/lambda", "github.com/aws/aws-sdk-go-v2/aws",
								"github.com/aws/aws-sdk-go-v2/config", "github.com/aws/aws-sdk-go-v2/credentials/stscreds",
								"github.com/aws/aws-sdk-go-v2/feature/s3/manager", "github.com/aws/aws-sdk-go-v2/service/ecr",
								// (continued list...)
							},
							Deny: []struct {
								Pkg  string `yaml:"pkg"`
								Desc string `yaml:"desc"`
							}{
								{Pkg: "log", Desc: "Use 'log \"github.com/sirupsen/logrus\"' instead"},
								{Pkg: "github.com/pkg/errors", Desc: "Should be replaced by standard lib errors package"},
							},
						},
					},
				},
			},
			Exclusions: Exclusions{
				Generated: "lax",
				Presets:   []string{"comments", "common-false-positives", "legacy", "std-error-handling"},
				Rules: []Rule{
					{Linters: []string{"revive"}, Path: "internal/pkg/data/trivy_types.go", Text: "exported: type name will be used as data.DataSource"},
					{Linters: []string{"funlen"}, Path: "internal/app/mcvs-scanner/data/data_integration_test.go", Text: "Function 'TestInsertImageSBOMIntoDB' is too long"},
				},
				Paths: []string{
					"internal/app/mcvs-scanner-cli/application/swagger", "third_party$", "builtin$", "examples$",
				},
			},
		},
		Formatters: FormatterSettings{
			Enable: []string{"gci", "gofmt", "gofumpt", "goimports"},
			Exclusions: struct {
				Generated string   `yaml:"generated"`
				Paths     []string `yaml:"paths"`
			}{
				Generated: "disable",
				Paths: []string{
					"internal/app/mcvs-scanner-cli/application/swagger", "third_party$", "builtin$", "examples$",
				},
			},
		},
	}

	// // Marshal the configuration to YAML
	// yamlData, err := yaml.Marshal(&config)
	// if err != nil {
	// 	log.Fatalf("error marshaling YAML: %v", err)
	// }

	// // Add the comment at the beginning and write to a file
	// yamlString := "# yamllint disable rule:line-length\n---\n" + string(yamlData)

	// // Write the YAML data to a file
	// filename := "config.yml" // Change this to your desired file name

	// err = os.WriteFile(filename, []byte(yamlString), 0o644)
	// if err != nil {
	// 	log.Fatalf("error writing to file: %v", err)
	// }
	filename := ".golangci.yml"
	// fmt.Printf("YAML content successfully written to %s\n", filename)

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	err := encoder.Encode(&config)
	if err != nil {
		fmt.Printf("Error encoding YAML: %v\n", err)
	}
	encoder.Close()
	yamlOutput := buffer.String()
	yamlOutput = strings.ReplaceAll(yamlOutput, "'", "\"")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
	}
	defer file.Close()
	_, err = file.WriteString("---\n" + yamlOutput)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func main() {
	golangFilePath := filepath.Join(".github", "workflows", "golang.yml")
	dependabotFilePath := filepath.Join(".github", "dependabot.yml")

	if err := createGolangWorkflow(golangFilePath); err != nil {
		fmt.Printf("Error writing Golang workflow file: %v\n", err)
	}

	if err := createDependabotConfig(dependabotFilePath); err != nil {
		fmt.Printf("Error writing Dependabot config file: %v\n", err)
	}

	bla()

	gitignore()

	golangCILint()

	fmt.Println("YAML files created successfully")
}
