package application

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/schubergphilis/mcvs-golang-action/internal/app/bla/data"
	"gopkg.in/yaml.v3"
)

func addStepComments(node *yaml.Node) {
	// Find jobs/mcvs-golang-action-taskfile-remote-url-ref-updater/steps/0
	var jobsNode *yaml.Node
	var updaterNode *yaml.Node
	var stepsSeq *yaml.Node

	// Find "jobs"
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == "jobs" {
			jobsNode = node.Content[i+1]
			break
		}
	}
	if jobsNode == nil {
		return
	}

	// Find "mcvs-golang-action-taskfile-remote-url-ref-updater"
	for i := 0; i < len(jobsNode.Content); i += 2 {
		if jobsNode.Content[i].Value == "mcvs-golang-action-taskfile-remote-url-ref-updater" {
			updaterNode = jobsNode.Content[i+1]
			break
		}
	}
	if updaterNode == nil {
		return
	}

	// Find "steps"
	for i := 0; i < len(updaterNode.Content); i += 2 {
		if updaterNode.Content[i].Value == "steps" {
			stepsSeq = updaterNode.Content[i+1]
			break
		}
	}
	if stepsSeq == nil || len(stepsSeq.Content) == 0 {
		return
	}

	// The first step node: add comments
	stepNode := stepsSeq.Content[0]
	stepNode.HeadComment = "yamllint disable rule:line-length"
	stepNode.FootComment = "yamllint enable rule:line-length"
}

func BuildWorkflow() data.Workflow {
	onNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{
				Kind:  yaml.ScalarNode,
				Value: "schedule",
			},
			{
				Kind: yaml.SequenceNode,
				Content: []*yaml.Node{
					{
						Kind: yaml.MappingNode,
						Content: []*yaml.Node{
							{
								Kind:  yaml.ScalarNode,
								Value: "cron",
							},
							{
								Kind:  yaml.ScalarNode,
								Value: "42 6 * * *",
							},
						},
					},
				},
			},
		},
	}

	steps := []*yaml.Node{
		{
			Kind:        yaml.MappingNode,
			HeadComment: "yamllint disable rule:line-length",
			Content: []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Value: "uses",
				},
				{
					Kind:  yaml.ScalarNode,
					Value: "schubergphilis/mcvs-golang-action-taskfile-remote-url-ref-updater@v0.1.2",
				},
			},
			FootComment: "yamllint enable rule:line-length",
		},
	}

	return data.Workflow{
		Name: "mcvs-golang-action-taskfile-remote-url-ref-updater",
		On:   onNode,
		Permissions: data.Permissions{
			Contents:     "write",
			PullRequests: "write",
		},
		Jobs: data.Jobs{
			McvsUpdater: data.Job{
				RunsOn: "ubuntu-24.04",
				Steps:  steps,
			},
		},
	}
}

func RenderAndSave(workflow data.Workflow, filePath string) error {
	outNode := &yaml.Node{}
	if err := outNode.Encode(workflow); err != nil {
		return err
	}
	for i := 0; i+1 < len(outNode.Content); i += 2 {
		k := outNode.Content[i]
		if k.Value == "on" {
			k.Style = yaml.DoubleQuotedStyle
		}
	}
	addStepComments(outNode)

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString("---\n"); err != nil {
		return err
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err = encoder.Encode(outNode)
	encoder.Close()
	if err != nil {
		return err
	}

	output := fixMisplacedSequenceComments(buf.String())
	_, err = f.WriteString(output)
	return err
}

func fixMisplacedSequenceComments(yaml string) string {
	lines := strings.Split(yaml, "\n")
	var result []string
	skipNext := false

	for i, line := range lines {
		// Skip this line if flagged
		if skipNext {
			skipNext = false
			continue
		}

		// If this is the misplaced comment at root level (0 indent)
		if line == "# yamllint enable rule:line-length" {
			// Remove preceding blank line if exists
			if len(result) > 0 && result[len(result)-1] == "" {
				result = result[:len(result)-1]
			}
			result = append(result, "      # yamllint enable rule:line-length")
			continue
		}

		// Also handle if it appears as a sequence item
		if matched, _ := regexp.MatchString(`^\s*-\s*#\s*yamllint enable rule:line-length`, line); matched {
			// Remove preceding blank line if exists
			if len(result) > 0 && result[len(result)-1] == "" {
				result = result[:len(result)-1]
			}
			result = append(result, "      # yamllint enable rule:line-length")
			continue
		}

		// Remove blank line after steps:
		if strings.TrimSpace(line) == "steps:" {
			result = append(result, line)
			if i+1 < len(lines) && lines[i+1] == "" {
				skipNext = true
			}
			continue
		}

		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func FindMcvsReposWithRemoteURLRef(rootDir string) ([]data.RepoInfo, error) {
	var repos []data.RepoInfo

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			repoPath := filepath.Join(rootDir, entry.Name())
			taskfilePath := filepath.Join(repoPath, "Taskfile.yml")
			if _, err := os.Stat(taskfilePath); err == nil {
				remoteURLRef, err := parseRemoteURLRefFromTaskfile(taskfilePath)
				if err != nil {
					fmt.Printf("Warning: Could not parse REMOTE_URL_REF in %s: %v\n", taskfilePath, err)
					continue
				}
				if remoteURLRef != "" {
					repos = append(repos, data.RepoInfo{
						Path:         repoPath,
						RemoteURLRef: remoteURLRef,
					})
				}
			}
		}
	}
	return repos, nil
}

func parseRemoteURLRefFromTaskfile(taskfile string) (string, error) {
	f, err := os.Open(taskfile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var taskfileYAML struct {
		Vars map[string]interface{} `yaml:"vars"`
	}
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&taskfileYAML); err != nil {
		return "", err
	}

	if val, ok := taskfileYAML.Vars["REMOTE_URL_REF"]; ok && val != nil {
		switch v := val.(type) {
		case string:
			return v, nil
		default:
			return fmt.Sprintf("%v", v), nil
		}
	}
	return "", nil
}

func hasChanges(repoPath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	// If output is empty, no changes
	return len(strings.TrimSpace(string(out))) > 0, nil
}

func UpdateWorkflowBranch(repoPath string) error {
	branchName := "update-workflows-" + time.Now().Format("20060102-150405")

	if err := runGit(repoPath, "checkout", "."); err != nil {
		return fmt.Errorf("git checkout . failed: %w", err)
	}

	if err := runGit(repoPath, "checkout", "main"); err != nil {
		return fmt.Errorf("git checkout main failed: %w", err)
	}

	if err := runGit(repoPath, "pull", "origin", "main"); err != nil {
		return fmt.Errorf("git pull origin main failed: %w", err)
	}

	if err := runGit(repoPath, "checkout", "-b", branchName); err != nil {
		return fmt.Errorf("git branch creation failed: %w", err)
	}

	filePath := ".github/workflows/mcvs-golang-action-taskfile-remote-url-ref-updater.yml"
	workflow := BuildWorkflow()
	if err := RenderAndSave(workflow, filepath.Join(repoPath, filePath)); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Workflow saved to", filePath)
	}

	changed, err := hasChanges(repoPath)
	if err != nil {
		return fmt.Errorf("failed to check git status: %w", err)
	}
	if !changed {
		fmt.Println("No changes detected, skipping commit and push")
		return nil
	}

	if err := runGit(repoPath, "add", "."); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	if err := runGit(repoPath, "commit", "-m", "build: update remote_url_ref in taskfile daily"); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	if err := runGit(repoPath, "push", "origin", branchName); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(out))
	}
	return nil
}
