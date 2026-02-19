# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

**mcvs-golang-action** is a GitHub composite action and reusable Taskfile that standardizes quality checks for Go projects. It provides:

- A GitHub Action ([action.yml](action.yml)) for CI/CD pipelines
- A remote Taskfile ([build/task.yml](build/task.yml)) for local development and CI automation

The action orchestrates: Go version installation (from go.mod), module verification, security scanning (osv-scanner, Grype, Trivy), golangci-lint, unit/integration/component tests, code coverage enforcement, and optional binary releases.

## Architecture

### Key Components

1. **action.yml** - Composite GitHub Action definition

   - Defines all inputs (build-tags, testing-type, release configs, etc.)
   - Orchestrates installation of Go and Task runner
   - Conditionally executes different testing/security/build workflows based on `testing-type` input
   - Supports multiple testing types: `unit`, `integration`, `component`, `coverage`, `lint`, `security-golang-modules`, `security-grype`, `security-trivy`, `graphql-lint`, `mcvs-texttidy`, `mocks-tidy`, `opa`

2. **build/task.yml** - Reusable Taskfile

   - Contains all task definitions that the action executes
   - Designed to be included remotely by other projects: `{{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/build/task.yml`
   - Defines versions for all tools (golangci-lint, osv-scanner, mockery, etc.)
   - Provides both CI-specific tasks (`test-cicd`, `coverage`) and development tasks (`test`, `lint`, `mocks`)

3. **scripts/package-version-updater.sh** - Automation
   - Periodically updates pinned tool versions in build/task.yml
   - Opens PRs with version updates using gh CLI

### Testing Architecture

Tests are organized using Go build tags:

- **`integration`** - Integration tests requiring external services (`*_integration_test.go` files with `//go:build integration`)
- **`component`** - Component tests (`*_component_test.go` files with `//go:build component`)
- **`e2e`** - End-to-end tests (`*_e2e_test.go` files with `//go:build e2e`)
- Default (no tag) - Unit tests

### Security Scanning Flow

1. **osv-scanner** - First line of defense for Go module vulnerabilities

   - Configured via `osv-scanner.toml` (see [osv-scanner.toml.example](osv-scanner.toml.example))
   - Allows temporary ignores (max 1 month) via `IgnoredVulns` with expiration dates
   - See [docs/osv-scanner.md](docs/osv-scanner.md) for detailed usage

2. **Grype** - Optional additional vulnerability scanning via Anchore

   - Triggered when `testing-type: security-grype`
   - Severity cutoff: HIGH or above

3. **Trivy** - Optional container/filesystem scanning
   - Triggered when `testing-type: security-trivy`
   - Supports custom ignore file via `.trivyignore`
   - Uses cached databases from public.ecr.aws to avoid rate limits

## Common Development Commands

### Using the Remote Taskfile (in consuming projects)

Set up `Taskfile.yml` in your project:

```yaml
version: 3
vars:
  REMOTE_URL: https://raw.githubusercontent.com
  REMOTE_URL_REF: v3.4.2 # Use latest stable version
  REMOTE_URL_REPO: schubergphilis/mcvs-golang-action
includes:
  remote: >-
    {{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/build/task.yml
```

Then run tasks with:

```bash
# Required: enable experimental remote taskfiles support
export TASK_X_REMOTE_TASKFILES=1

# Run unit tests
task remote:test --yes

# Run integration tests (includes unit tests)
task remote:test-integration --yes

# Run component tests
task remote:test-component --yes

# Run linting
task remote:lint --yes

# Run with custom build tags
BUILD_TAGS="integration" task remote:lint --yes

# Run code coverage
task remote:coverage --yes

# Run security scanning
task remote:osv-scanner --yes

# Automatically fix linting issues
task remote:fix-linting-issues --yes

# List all available tasks
task --list-all
```

### Fixing Linting Issues

The `fix-linting-issues` task automatically fixes common linting problems:

```bash
task remote:fix-linting-issues --yes
```

This task uses:
- **golines** (v0.12.2) - Reformats code to meet line length requirements by intelligently wrapping long lines
- **wsl** (v5.1.0) - Fixes whitespace linting issues by adding/removing blank lines according to Go style guidelines

The task runs:
1. `golines . -w` - Reformats all Go files in the current directory
2. `wsl -fix ./...` - Fixes whitespace issues in all Go packages

After running, review the changes as some linting issues may still require manual intervention.

### Testing in This Repository

This repository uses itself for CI. Check [.github/workflows/golang.yml](.github/workflows/golang.yml) for examples:

```bash
# The workflows use the action locally with uses: ./
# Matrix strategy runs multiple testing-type values in parallel
```

### Overriding Variables

Override Taskfile variables when including remotely:

```yaml
includes:
  remote:
    taskfile: >-
      {{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/build/task.yml
    vars:
      CODE_COVERAGE_STRICT: "false" # Disable strict coverage enforcement
      GOLANGCI_LINT_RUN_TIMEOUT_MINUTES: "5" # Increase timeout
```

Available override variables (see [build/task.yml](build/task.yml) lines 47-53):

- `CODE_COVERAGE_STRICT` - Enforce minimum coverage (default: "true")
- `GOLANGCI_LINT_CONFIG_PATH` - Path to golangci-lint config (default: ".golangci.yml")
- `GOLANGCI_LINT_RUN_TIMEOUT_MINUTES` - Linter timeout (default: 3)
- `BUILD_TAGS` - Build tags for tests/linting (default: "component,e2e,integration")

## Using the GitHub Action

### Basic Usage

```yaml
name: Golang
on: pull_request
permissions:
  contents: read
  packages: read
jobs:
  MCVS-golang-action:
    strategy:
      matrix:
        args:
          - testing-type: "unit"
          - testing-type: "lint"
          - testing-type: "coverage"
          - testing-type: "security-golang-modules"
    runs-on: ubuntu-24.04
    env:
      TASK_X_REMOTE_TASKFILES: 1
    steps:
      - uses: actions/checkout@v4.1.1
      - uses: schubergphilis/mcvs-golang-action@v3 # Use @v3 for latest v3.x.x
        with:
          testing-type: ${{ matrix.args.testing-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

### Advanced Usage with Releases

For building binaries on tagged releases:

```yaml
- uses: schubergphilis/mcvs-golang-action@v3
  with:
    release-application-name: "my-app"
    release-architecture: "amd64"
    release-dir: "./cmd/my-app"
    release-os: "linux"
    release-type: "binary"
    release-build-tags: "lambda.norpc" # Optional, for AWS Lambda builds
    token: ${{ secrets.GITHUB_TOKEN }}
```

### Key Action Inputs

- **testing-type** - Main selector: `unit`, `integration`, `component`, `coverage`, `lint`, `security-golang-modules`, `security-grype`, `security-trivy`
- **build-tags** - Build constraints for tests/linting (e.g., "integration,component")
- **golang-unit-tests-exclusions** - Regex to exclude packages from unit tests (e.g., `\(cmd\/app\|internal\/app\)`)
- **code-coverage-expected** - Minimum coverage percentage (default: 80)
- **golangci-timeout** - Timeout for golangci-lint
- **test-timeout** / **code-coverage-timeout** - Timeouts for test execution (e.g., "10m0s")

## Important Implementation Details

### Module Verification

- `go mod tidy` is automatically run and verified (fails if git diff shows changes)
- Private Go modules: use `github-token-for-downloading-private-go-modules` input to configure git credentials

### Linting with Multiple Build Tags

When you have test files with different build tags, lint them separately:

```yaml
- testing-type: "lint"  # Lint main code
- testing-type: "lint", build-tags: "integration"
- testing-type: "lint", build-tags: "component"
```

### Tool Versions

All tool versions are pinned in [build/task.yml](build/task.yml) (lines 52-109):

- golangci-lint: v2.9.0
- osv-scanner: v2.3.3
- mockery: v3.6.4
- opa/regal: v1.13.1 / v0.38.1
- golines: v0.12.2
- wsl: v5.1.0
- Task runner: 3.46.4 (defined in action.yml)

The [package-version-updater workflow](.github/workflows/package-version-updater.yml) automatically opens PRs to update these versions weekly.

### Go Version Management

The Go version is determined from `go.mod` using `actions/setup-go` with `go-version-file: go.mod`. This ensures CI uses the same Go version as defined in the project.

## Project-Specific Notes

### This Repository Structure

```
.
├── action.yml              # Composite action definition
├── build/
│   └── task.yml           # Remote Taskfile with all tasks
├── scripts/
│   └── package-version-updater.sh  # Tool version updater
├── .github/workflows/
│   ├── golang.yml         # Self-testing workflow (uses this action)
│   └── package-version-updater.yml  # Weekly tool updates
├── docs/
│   ├── osv-scanner.md     # OSV scanner usage guide
│   └── presentations/     # Present tool presentations
├── osv-scanner.toml.example  # Example vulnerability ignore config
└── go.mod                 # Module definition (no .go files in this repo)
```

### No Application Code

This repository contains no `.go` application code - it's purely tooling. The `go.mod` exists to install tools via `go install` and to define the Go version for the action.

### Versioning

- Use `@v3` in workflows to automatically get latest v3.x.x updates
- Use specific tags (e.g., `@v3.4.2`) for pinned versions
- Breaking changes only occur on major version bumps (v3 → v4)
- Check the [releases page](https://github.com/schubergphilis/mcvs-golang-action/releases) for changelog
