# MCVS-golang-action

The Mission Critical Vulnerability Scanner (MCVS) Golang Action repository is a
collection of standardized tools to ensure a certain level of quality of a
project with Go code.

## Github Action

The [GitHub Action](https://github.com/features/actions) in this repository
consists of the following steps:

- Install the Golang version that is defined in the project `go.mod`.
- Verify to be downloaded Golang modules.
- Check for incorrect import order and indicate how to resolve it.
- Code security scanning and suppression of certain CVEs for a maximum of one
  month. In some situations a particular CVE will be resolved in a couple of
  weeks and this allows the developer to continue in a safe way while knowing
  that the pipeline will fail again if the issue has not been resolved in a
  couple of weeks.
- Linting.
- Unit tests.
- Integration tests.
- Code coverage.

In summary, using this action will ensure that Golang code meets certain
standards before it will be deployed to production as the assembly line will
fail if an issue arises.

Note: there is an [internal action](.github/workflows/package-version-updater.yml)
that will update package versions that cannot be updated by Dependabot.

## Versioning

This action follows semantic versioning. When using this action in your workflows:

- **Latest stable version**: Use the latest `v3.x.x` tag (e.g., `v3.4.2`) for production workflows
- **Major version tracking**: Use `@v3` to automatically get the latest v3.x.x updates
- **Taskfile references**: When including the remote Taskfile, use a specific version tag (e.g., `v2.1.0`) that matches your needs
- **Breaking changes**: Major version bumps (v3 → v4) may introduce breaking changes and require workflow updates

Check the [releases page](https://github.com/schubergphilis/mcvs-golang-action/releases) for the latest version and changelog.

## Taskfile

Another tool is configuration for [Task](https://taskfile.dev/). This repository
offers a `./build/task.yml` which contains standard tasks, like installing and
running a linter.

This `./build/task.yml` can then be used by other projects. This has the
advantage that you do not need to copy and paste Makefile snippets from one
project to another. As a consequence each project using this `./build/task.yml`
immediately benefits from improvements made here (e.g. new tasks or
improvements in the tasks).

If you are new to Task, you may want to check out the following resources:

- [Installation instructions](https://taskfile.dev/installation/)
- Instructions to [configure completions](https://taskfile.dev/installation/#setup-completions)
- [Integrations](https://taskfile.dev/integrations/) with e.g. Visual Studio Code, Sublime and IntelliJ.

### Configuration

The `./build/task.yml` in this project defines a number of variables. Some of
these can be overridden when including this Taskfile in your project. See the
example below, where the `CODE_COVERAGE_STRICT` variable is overridden, for how
to do this.

The following variables can be overridden:

| Variable                    | Description                                                                                              |
| :-------------------------- | :------------------------------------------------------------------------------------------------------- |
| `CODE_COVERAGE_STRICT`      | Enables or disables strict enforcement of setting the minimum coverage to the maximum observed coverage. |
| `GOLANGCI_LINT_CONFIG_PATH` | Defines the path to the golangci-lint configuration file.                                                |

## Usage

### Locally

Create a `Taskfile.yml` with the following content:

```yml
---
version: 3

vars:
  REMOTE_URL: https://raw.githubusercontent.com
  REMOTE_URL_REF: v2.1.0
  REMOTE_URL_REPO: schubergphilis/mcvs-golang-action

includes:
  remote: >-
    {{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/build/task.yml
```

and run:

```zsh
TASK_X_REMOTE_TASKFILES=1 \
task remote:test
```

Note that the `TASK_X_REMOTE_TASKFILES` variable is required as long as the
remote Taskfiles are still experimental. (See [issue
1317](https://github.com/go-task/task/issues/1317) for more information.)

You can use `task --list-all` to get a list of all available tasks.
Alternatively, if you have [configured
completions](https://taskfile.dev/installation/#setup-completions) in your
shell, you can tab to get a list of available tasks.

If you want to override one of the variables in our Taskfile, you will have to
adjust the `includes` sections like this:

```yml
---
includes:
  remote:
    taskfile: >-
      {{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/build/task.yml
    vars:
      CODE_COVERAGE_STRICT: "false"
```

Note: same goes for the `GOLANGCI_LINT_RUN_TIMEOUT_MINUTES` setting.

## Build Tags

Build tags (also known as build constraints) allow you to include or exclude Go files from compilation based on conditions. This action supports the following common build tag patterns:

- **`integration`**: For integration tests that require external services or databases
- **`component`**: For component tests that test multiple units working together
- **`e2e`**: For end-to-end tests that test the entire application flow
- **`lambda.norpc`**: For building AWS Lambda functions without RPC support

### Using Build Tags

When running tests with specific build tags:
```zsh
# Run integration tests
task remote:test-integration --yes

# Run component tests
task remote:test-component --yes
```

When linting code with specific build tags, you may need to run the linter multiple times to cover all code paths:
```yml
- testing-type: "lint"  # Lint main code
- testing-type: "lint", build-tags: "integration"  # Lint integration test code
- testing-type: "lint", build-tags: "component"  # Lint component test code
```

This ensures that code in test files with different build tags is properly linted.

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
    env:
      TASK_X_REMOTE_TASKFILES: 1
    steps:
      - uses: actions/checkout@v4.1.1
      - uses: schubergphilis/mcvs-golang-action@v3
        with:
          testing-type: ${{ matrix.args.testing-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

This basic configuration will run unit tests, linting, code coverage checks, and security scanning on your Go code.

#### Advanced Example

For projects with multiple build configurations, integration tests, or custom requirements, create a `.github/workflows/golang.yml` file with the following content:

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
          - release-architecture: "amd64",
            release-dir: "./cmd/path-to-app",
            release-type: "binary",
            release-application-name: "some-app",
          - release-architecture: "arm64",
            release-dir: "./cmd/path-to-app",
            release-type: "binary",
            release-application-name: "some-lambda-func",
            release-build-tags: "lambda.norpc",
          - testing-type: "component"
          - testing-type: "coverage"
          - testing-type: "graphql-lint"
          - testing-type: "integration"
          - testing-type: "lint", build-tags: "component"
          - testing-type: "lint", build-tags: "e2e"
          - testing-type: "lint", build-tags: "integration"
          - testing-type: "mcvs-texttidy"
          - testing-type: "mocks-tidy"
          - testing-type: "security-golang-modules"
          - testing-type: "security-grype"
          - testing-type: "security-trivy"
            security-trivyignore: ""
          - testing-type: "unit"
    runs-on: ubuntu-24.04
    env:
      TASK_X_REMOTE_TASKFILES: 1
      test-timeout: 10m0s
    steps:
      - uses: actions/checkout@v4.1.1
        with:
          fetch-depth: 0 # this is necessary for gta partial testing
      - uses: schubergphilis/mcvs-golang-action@v0.9.0
        with:
          build-tags: ${{ matrix.args.build-tags }}
          golang-unit-tests-exclusions: |-
            \(cmd\/some-app\|internal\/app\/some-app\)
          gta-base-branch: main
          gta-partial-testing: true
          release-architecture: ${{ matrix.args.release-architecture }}
          release-dir: ${{ matrix.args.release-dir }}
          release-type: ${{ matrix.args.release-type }}
          security-trivyignore: ${{ matrix.args.security-trivyignore }}
          task-install: yes
          testing-type: ${{ matrix.args.testing-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
          test-timeout: ${{ env.test-timeout }}
          code-coverage-timeout: ${{ env.test-timeout }}
```

and a [.golangci.yml](https://golangci-lint.run/usage/configuration/).

<!-- markdownlint-disable MD013 -->

| Option                                          | Default | Required | Description                                                                                                      |
| :---------------------------------------------- | :------ | -------- | :--------------------------------------------------------------------------------------------------------------- |
| build-tags                                      | x       |          | Build tags to use when running tests and linting (e.g., "integration", "component", "e2e")                       |
| code-coverage-expected                          | x       |          | Minimum expected code coverage percentage for standard tests                                                     |
| code-coverage-opa-expected                      | x       |          | Minimum expected code coverage percentage for OPA (Open Policy Agent) tests                                      |
| code-coverage-timeout                           |         |          | Timeout duration for code coverage analysis (e.g., "10m0s")                                                      |
| github-token-for-downloading-private-go-modules |         |          | GitHub token with permissions to download Go modules from private repositories                                   |
| golangci-timeout                                | x       |          | Timeout duration for golangci-lint execution                                                                     |
| golang-unit-tests-exclusions                    | x       |          | Regex pattern to exclude specific packages from unit testing (e.g., `\(cmd\/app\|internal\/app\)`)               |
| grype-version                                   |         |          | Specific version of Grype vulnerability scanner to use                                                           |
| gta-base-branch                                 | x       |          | The branch changed go packages will be compared to, to perform partial tests                                     |           
| gta-partial-testing                             | x       |          | Whether to run partial tests (true or false)                                                                     |
| release-application-name                        |         |          | Name of the application binary to build (required when release-type is set)                                      |
| release-architecture                            |         |          | Target architecture for the binary (e.g., "amd64", "arm64")                                                      |
| release-build-tags                              |         |          | Build tags to use when building the release binary (e.g., "lambda.norpc")                                        |
| release-dir                                     |         |          | Directory containing the main.go file for the binary to build                                                    |
| release-os                                      | x       |          | Target operating system for the binary (e.g., "linux", "darwin")                                                 |
| release-type                                    |         |          | Type of release to build (e.g., "binary")                                                                        |
| task-install                                    | x       |          | Whether to install Task runner ("yes" or "no")                                                                   |
| task-version                                    | x       |          | Version of Task runner to install                                                                                |
| testing-type                                    |         |          | Type of testing to run (e.g., "unit", "integration", "lint", "coverage", "security-golang-modules")              |
| test-timeout                                    |         |          | Timeout duration for test execution (e.g., "10m0s")                                                              |
| token                                           |         |          | GitHub token for authentication (typically ${{ secrets.GITHUB_TOKEN }})                                          |
| trivy-action-db                                 | x       |          | Trivy vulnerability database configuration                                                                       |
| trivy-action-java-db                            | x       |          | Trivy Java vulnerability database configuration                                                                  |


Note: If an **x** is registered in the Default column, refer to the
[action.yml](action.yml) for the corresponding value.

<!-- markdownlint-enable MD013 -->

### Releases

In some cases, you may want the executable binary to be built and released
automatically. This action will build the binary which could then be used
as a release asset.

Create a `.github/workflows/golang-releases.yml` file with the following
content:

```yml
---
name: golang-releases
"on": push
permissions:
  contents: write
  packages: read
jobs:
  mcvs-golang-action:
    strategy:
      matrix:
        args:
          - release-application-name: mcvs-image-downloader
            release-architecture: amd64
            release-dir: cmd/mcvs-image-downloader
            release-type: binary
          - release-application-name: mcvs-image-downloader
            release-architecture: arm64
            release-dir: cmd/mcvs-image-downloader
            release-os: darwin
            release-type: binary
    runs-on: ubuntu-24.04
    env:
      TASK_X_REMOTE_TASKFILES: 1
    steps:
      - uses: actions/checkout@v4.2.2
      - uses: schubergphilis/mcvs-golang-action@v3.4.2
        with:
          release-application-name: ${{ matrix.args.release-application-name }}
          release-architecture: ${{ matrix.args.release-architecture }}
          release-build-tags: ${{ matrix.args.release-build-tags }}
          release-dir: ${{ matrix.args.release-dir }}
          release-os: ${{ matrix.args.release-os }}
          release-type: ${{ matrix.args.release-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

### Integration

To execute integration tests, make sure that the code is located in a file with
a `_integration_test.go` postfix, such as `some_integration_test.go`.
Additionally, include the following header in the file:

```bash
//go:build integration
```

After adding this header, issue the command `task remote:test-integration --yes`
as demonstrated in this example. This action will run both unit and integration
tests. If `task remote:test --yes` is executed, only unit tests will be run.

### Component

See the integration paragraph for the steps and replace `integration` with
`component` to run them.

### Downloading released assets from another private repository

You will need a personal access token (PAT) with the `repo` scope. To download
releases from a private repository. You can simply use the gh command or curl
to download the release assets. Please read the
[GitHub documentation](https://docs.github.com/en/rest/releases/assets)
for more information.
