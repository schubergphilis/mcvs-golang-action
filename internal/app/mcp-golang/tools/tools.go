package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// Handler is a struct that handles tool requests.
type Handler struct{}

// NewHandler creates a new Handler instance.
func NewHandler() *Handler {
	return &Handler{}
}

func (t *Handler) HandleLint(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := []string{
		"remote:lint",
		"-y",
	}

	dir, ok := request.Params.Arguments["directory"].(string)
	if ok {
		arguments = append(arguments, "-d", dir)
	}
	return t.executeTask(arguments)
}

func (t *Handler) HandleUnitTest(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := []string{
		"remote:test",
		"-y",
	}

	dir, ok := request.Params.Arguments["directory"].(string)
	if ok {
		arguments = append(arguments, "-d", dir)
	}
	extraTestArguments := []string{}
	verbose, ok := request.Params.Arguments["verbose"].(bool)
	if ok && verbose {
		extraTestArguments = append(extraTestArguments, "-v")
	}
	coverage, ok := request.Params.Arguments["coverage"].(bool)
	if ok && coverage {
		extraTestArguments = append(extraTestArguments, "-cover")
	}
	testCase, ok := request.Params.Arguments["test-case"].(string)
	if ok && testCase != "" {
		extraTestArguments = append(extraTestArguments, "-run", testCase)
	}
	if len(extraTestArguments) > 0 {
		arguments = append(arguments, fmt.Sprintf(`TEST_EXTRA_ARGS=%s`, strings.Join(extraTestArguments, " ")))
	}
	buildTags, ok := request.Params.Arguments["build-tags"].(string)
	if ok && buildTags != "none" {
		arguments = append(arguments, fmt.Sprintf("TEST_TAGS=%s", buildTags))
	}
	return t.executeTask(arguments)
}

type memWriter struct {
	bytes []byte
}

func (m *memWriter) Write(p []byte) (n int, err error) {
	m.bytes = append(m.bytes, p...)
	return len(p), nil
}

func (m *memWriter) String() string {
	return string(m.bytes)
}

func (t *Handler) executeTask(arguments []string) (*mcp.CallToolResult, error) {
	cmd := exec.Command("task", arguments...)
	stdErr := &memWriter{}
	stdOut := &memWriter{}
	cmd.Stderr = stdErr
	cmd.Stdout = stdOut
	err := cmd.Run()
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("%s %s %s", stdErr.String(), stdOut.String(), err)), nil
	}

	return mcp.NewToolResultText(stdOut.String()), nil
}
