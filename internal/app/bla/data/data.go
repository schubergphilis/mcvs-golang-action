package data

import "gopkg.in/yaml.v3"

type Workflow struct {
	Name        string      `yaml:"name"`
	On          *yaml.Node  `yaml:"on"`
	Permissions Permissions `yaml:"permissions"`
	Jobs        Jobs        `yaml:"jobs"`
}

type Permissions struct {
	Contents     string `yaml:"contents"`
	PullRequests string `yaml:"pull-requests"`
}

type Jobs struct {
	McvsUpdater Job `yaml:"mcvs-golang-action-taskfile-remote-url-ref-updater"`
}

type Job struct {
	RunsOn string       `yaml:"runs-on"`
	Steps  []*yaml.Node `yaml:"steps"`
}

type StepUses struct {
	Uses string `yaml:"uses"`
}

type RepoInfo struct {
	Path         string
	RemoteURLRef string
}
