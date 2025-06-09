package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/schubergphilis/mcvs-golang-action/internal/app/mcp-golang/tools"
)

func main() {
	s := server.NewMCPServer(
		"MCVS Golang Action MCP Server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	toolHandler := tools.NewHandler()

	s.AddTool(
		mcp.NewTool(
			"static-analysis",
			mcp.WithDescription("Tool to lint the project (runs static analysis on the codebase) leveraging golangci-lint"),

			mcp.WithString("directory",
				mcp.Description("The directory to run the static analysis in"),
			),
		),
		toolHandler.HandleLint,
	)
	s.AddTool(
		mcp.NewTool(
			"static-analysis-with-json-output",
			mcp.WithDescription("Tool to lint the project (runs static analysis on the codebase) leveraging golangci-lint, but with JSON output"),

			mcp.WithString("directory",
				mcp.Description("The directory to run the static analysis in"),
			),
		),
		toolHandler.HandleLintWithJSONOutput,
	)
	s.AddTool(
		mcp.NewTool("test",
			mcp.WithDescription("Tool to run tests on the project"),

			// We want to ensure the directory is always set since
			// copilot will not execute mcp server in the project root.
			mcp.WithString("directory",
				mcp.Description("The directory to run the tests in"),
				mcp.Required(),
			),
			mcp.WithBoolean("verbose",
				mcp.Description("Enable verbose output"),
				mcp.DefaultBool(false),
			),
			mcp.WithString("test-case",
				mcp.Description("The test case to run"),
			),
			mcp.WithBoolean("coverage",
				mcp.Description("Enable coverage report"),
				mcp.DefaultBool(false),
			),
			mcp.WithString("build-tags",
				mcp.Description("Build tags to use when running tests"),
				mcp.Enum("none", "integration", "component", "e2e"),
				mcp.DefaultString("none"),
			),
		),
		toolHandler.HandleUnitTest,
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
