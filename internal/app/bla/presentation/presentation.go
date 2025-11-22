package presentation

import (
	"fmt"

	"github.com/schubergphilis/mcvs-golang-action/internal/app/bla/application"
)

func RunCLI() {
	rootDir := "../"
	repos, err := application.FindMcvsReposWithRemoteURLRef(rootDir)
	if err != nil {
		fmt.Println("Error scanning repos:", err)
		return
	}
	for _, repo := range repos {
		fmt.Printf("Repo: %s\n  REMOTE_URL_REF: %s\n", repo.Path, repo.RemoteURLRef)

		fmt.Printf("Processing repo: %s\n", repo.Path)
		err := application.UpdateWorkflowBranch(repo.Path)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Println("  Updated and pushed branch for:", repo.Path)
		}
	}
}
