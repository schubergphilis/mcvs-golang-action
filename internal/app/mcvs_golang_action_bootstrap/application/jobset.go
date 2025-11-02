package application

import "fmt"

const (
	Version = "v3.7.15"

	DirWorkflows = ".github/workflows"
	URLRawBase   = "https://raw.githubusercontent.com/schubergphilis/mcvs-golang-action/refs/tags"

	URLGeneralWorkflow = URLRawBase + "/%s/" + DirWorkflows + "/general.yml"
	URLGoModUpdater    = URLRawBase + "/%s/" + DirWorkflows + "/gomod-go-version-updater.yml"
	URLDependabot      = URLRawBase + "/%s/.github/dependabot.yml"

	PathGeneralWorkflow = DirWorkflows + "/general.yml"
	PathGoModUpdater    = DirWorkflows + "/gomod-go-version-updater.yml"
	PathDependabot      = ".github/dependabot.yml"
)

type DownloadJob struct {
	URL         string
	Destination string
}

// HardcodedJobs returns the unmodifiable list of jobs to download.
func HardcodedJobs() []DownloadJob {
	return []DownloadJob{
		{
			URL:         fmt.Sprintf(URLGeneralWorkflow, Version),
			Destination: PathGeneralWorkflow,
		},
		{
			URL:         fmt.Sprintf(URLDependabot, Version),
			Destination: PathDependabot,
		},
		{
			URL:         fmt.Sprintf(URLGoModUpdater, Version),
			Destination: PathGoModUpdater,
		},
	}
}
