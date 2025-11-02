package presentation

import (
	"log/slog"
	"os"

	"github.com/schubergphilis/mcvs-golang-action/internal/app/mcvs_golang_action_bootstrap/application"
)

type Presenter interface {
	Run() error
}

type CLIPresenter struct {
	Executor application.Executor
	Logger   *slog.Logger
}

func getModuleName() (string, error) {
	path := "go.mod"

	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	for _, line := range splitLines(string(b)) {
		if len(line) > 7 && line[:7] == "module " {
			return line[7:], nil
		}
	}

	return "", os.ErrNotExist
}

func splitLines(s string) []string {
	var lines []string

	l := 0

	for i := range s {
		if s[i] == '\n' {
			lines = append(lines, s[l:i])
			l = i + 1
		}
	}

	if l < len(s) {
		lines = append(lines, s[l:])
	}

	return lines
}

func (p *CLIPresenter) Run() error {
	jobs := application.HardcodedJobs()
	p.Logger.Info("Starting batch download", "count", len(jobs))

	if err := p.Executor.DownloadFiles(jobs); err != nil {
		p.Logger.Error("Batch download failed", "err", err)

		return err
	}

	p.Logger.Info("All downloads successful")

	p.Logger.Info("Scanning and writing .mockery.yaml")

	modName, err := getModuleName()
	if err != nil {
		p.Logger.Error("Could not determine go module name", "err", err)

		return err
	}

	cwd, _ := os.Getwd()

	err = p.Executor.GenerateMockeryYaml(modName, cwd, ".mockery.yaml", p.Logger)
	if err != nil {
		p.Logger.Error("Failed to write .mockery.yaml", "err", err)

		return err
	}

	p.Logger.Info(".mockery.yaml written successfully")

	return nil
}
