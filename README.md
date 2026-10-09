# MCVS Golang Action

[![GitHub release](https://img.shields.io/github/v/release/schubergphilis/mcvs-golang-action)](https://github.com/schubergphilis/mcvs-golang-action/releases)
[![License](https://img.shields.io/github/license/schubergphilis/mcvs-golang-action)](LICENSE)

<img src="./assets/logos/mcvs-golang-action.png" width="250"></a>

The Mission Critical Vulnerability Scanner (MCVS) Golang Action repository is a
collection of standardized tools to ensure a certain level of quality of a
project with Go code.

## Quickstart

### Locally

Create a `Taskfile.yml` with the following content:

```yml
---
version: 3

vars:
  REMOTE_URL: https://raw.githubusercontent.com
  REMOTE_URL_REF: v3.12.12
  REMOTE_URL_REPO: schubergphilis/mcvs-golang-action

includes:
  remote: >-
    {{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/build/task.yml
```

and run one of the most common tasks:

```zsh
task remote:test --yes
task remote:test-integration --yes
task remote:lint --yes
task remote:coverage --yes
task remote:osv-scanner --yes
```

You can use `task --list-all` to get a list of all available tasks.
Alternatively, if you have [configured
completions](https://taskfile.dev/installation/#setup-completions) in your
shell, you can tab to get a list of available tasks.

### GitHub

#### Basic Example

For a simple project that needs standard testing and linting, create a `.github/workflows/golang.yml` file:

```yml
---
name: Golang
"on": pull_request
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
    steps:
      - uses: actions/checkout@v4.1.1
      - uses: schubergphilis/mcvs-golang-action@v3
        with:
          testing-type: ${{ matrix.args.testing-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

This basic configuration will run unit tests, linting, code coverage checks, and security scanning on your Go code.

## Documentation

- [GitHub Action](docs/overview.md): the steps the action runs.
- [Versioning](docs/versioning.md): which version tag to use.
- [Taskfile](docs/taskfile.md): what the remote Taskfile offers and how to
  get started with Task.
- [Action inputs](docs/action-inputs.md): advanced workflow example and all
  inputs of the GitHub Action.
- [Taskfile variables](docs/taskfile-variables.md): variables that can be
  overridden when including the Taskfile.
- [Build tags](docs/build-tags.md): integration, component and e2e tests.
- [Linting](docs/linting.md): automatically fixing linting issues.
- [Releases](docs/releases.md): building binaries as release assets.
- [osv-scanner](docs/osv-scanner.md): security scanning and temporarily
  ignoring vulnerabilities.
