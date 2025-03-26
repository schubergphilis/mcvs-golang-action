package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

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
			"MCVS-golang-action": {
				Strategy: Strategy{
					Matrix: Matrix{
						Args: []map[string]interface{}{
							{
								"release-application-name": "some-app",
								"release-architecture":     "amd64",
								"release-dir":              "./cmd/path-to-app",
								"release-type":             "binary",
							},
							{
								"release-application-name": "some-lambda-func",
								"release-architecture":     "arm64",
								"release-build-tags":       "lambda.norpc",
								"release-dir":              "./cmd/path-to-app",
								"release-type":             "binary",
							},
							{"testing-type": "component"},
							{"testing-type": "coverage"},
							{"testing-type": "integration"},
							{"testing-type": "lint", "build-tags": "component"},
							{"testing-type": "lint", "build-tags": "e2e"},
							{"testing-type": "lint", "build-tags": "integration"},
							{"testing-type": "mcvs-texttidy"},
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
						Uses: "schubergphilis/mcvs-golang-action@v0.9.0",
						With: map[string]any{
							"build-tags":                   "${{ matrix.args.build-tags }}",
							"golang-unit-tests-exclusions": golangUnitTestsExclusions,
							"release-architecture":         "${{ matrix.args.release-architecture }}",
							"release-dir":                  "${{ matrix.args.release-dir }}",
							"release-type":                 "${{ matrix.args.release-type }}",
							"security-trivyignore":         "${{ matrix.args.security-trivyignore }}",
							"testing-type":                 "${{ matrix.args.testing-type }}",
							"token":                        "${{ secrets.GITHUB_TOKEN }}",
							"test-timeout":                 "${{ env.test-timeout }}",
							"code-coverage-timeout":        "${{ env.test-timeout }}",
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

func main() {
	golangFilePath := filepath.Join(".github", "workflows", "golang.yml")
	dependabotFilePath := filepath.Join(".github", "dependabot.yml")

	if err := createGolangWorkflow(golangFilePath); err != nil {
		fmt.Printf("Error writing Golang workflow file: %v\n", err)
	}

	if err := createDependabotConfig(dependabotFilePath); err != nil {
		fmt.Printf("Error writing Dependabot config file: %v\n", err)
	}

	fmt.Println("YAML files created successfully")
}
