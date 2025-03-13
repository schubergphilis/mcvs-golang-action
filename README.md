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

## Taskfile

Another tool is configuration for [Task](https://taskfile.dev/). This repository
offers a `Taskfile.yml` which contains standard tasks, like installing and
running a linter.

This `Taskfile.yml` can then be used by other projects. This has the advantage
that you do not need to copy and paste Makefile snippets from one project to
another. As a consequence each project using this `Taskfile.yml` immediately
benefits from improvements made here (e.g. new tasks or improvements in the
tasks).

If you are new to Task, you may want to check out the following resources:

- [Installation instructions](https://taskfile.dev/installation/)
- Instructions to [configure completions](https://taskfile.dev/installation/#setup-completions)
- [Integrations](https://taskfile.dev/integrations/) with e.g. Visual Studio Code, Sublime and IntelliJ.

### Configuration

The `Taskfile.yml` in this project defines a number of variables. Some of these
can be overridden when including this Taskfile in your project. See the example
below, where the `GCI_SECTIONS` variable is overridden, for how to do this.

The following variables can be overridden:

| Variable      | Description                                                                                                                     |
| :------------ | :------------------------------------------------------------------------------------------------------------------------------ |
| `GCI_SECTION` | Define how `gci` processes inputs (see the [gci README](https://github.com/daixiang0/gci?tab=readme-ov-file#usage) for details) |

## Usage

### Locally

Create a `Taskfile.yml` with the following content:

```yml
---
version: 3

vars:
  REMOTE_URL: https://raw.githubusercontent.com
  REMOTE_URL_REF: v0.10.2
  REMOTE_URL_REPO: schubergphilis/mcvs-golang-action

includes:
  remote: >-
    {{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/Taskfile.yml
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
      {{.REMOTE_URL}}/{{.REMOTE_URL_REPO}}/{{.REMOTE_URL_REF}}/Taskfile.yml
    vars:
      GCI_SECTIONS: >-
        -s standard
        -s default
        -s alias
```

Note: same goes for the `GOLANGCI_LINT_RUN_TIMEOUT_MINUTES` setting.

### GitHub

Create a `.github/workflows/golang.yml` file with the following content:

```yml
---
name: Golang
"on": push
permissions:
  contents: read
  packages: read
jobs:
  MCVS-golang-action:
    strategy:
      matrix:
        args:
          - {
              release-architecture: 'amd64',
              release-dir: './cmd/path-to-app',
              release-type: 'binary',
              release-application-name: 'some-app',
            }
          - {
              release-architecture: 'arm64',
              release-dir: './cmd/path-to-app',
              release-type: 'binary',
              release-application-name: 'some-lambda-func',
              release-build-tags: 'lambda.norpc',
            }
          - { testing-type: 'component' }
          - { testing-type: 'coverage' }
          - { testing-type: 'integration' }
          - { testing-type: 'lint', build-tags: 'component' }
          - { testing-type: 'lint', build-tags: 'e2e' }
          - { testing-type: 'lint', build-tags: 'integration' }
          - { testing-type: 'mcvs-texttidy' }
          - { testing-type: 'security-golang-modules' }
          - { testing-type: 'security-grype' }
          - { testing-type: 'security-trivy', security-trivyignore: '' }
          - { testing-type: 'unit' }
    runs-on: ubuntu-22.04
    env:
      TASK_X_REMOTE_TASKFILES: 1
      test-timeout: 10m0s
    steps:
      - uses: actions/checkout@v4.1.1
      - uses: schubergphilis/mcvs-golang-action@v0.9.0
        with:
          build-tags: ${{ matrix.args.build-tags }}
          golang-unit-tests-exclusions: |-
            \(cmd\/some-app\|internal\/app\/some-app\)
          release-architecture: ${{ matrix.args.release-architecture }}
          release-dir: ${{ matrix.args.release-dir }}
          release-type: ${{ matrix.args.release-type }}
          security-trivyignore: ${{ matrix.args.security-trivyignore }}
          testing-type: ${{ matrix.args.testing-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
          test-timeout:  ${{ env.test-timeout }}
          code-coverage-timeout: ${{ env.test-timeout }}
```

and a [.golangci.yml](https://golangci-lint.run/usage/configuration/).

<!-- markdownlint-disable MD013 -->

| Option                                          | Default | Required |
| :---------------------------------------------- | :------ | -------- |
| build-tags                                      | x       |          |
| code-coverage-expected                          | x       |          |
| code-coverage-timeout                           |         |          |
| gci                                             | x       |          |
| github-token-for-downloading-private-go-modules |         |          |
| golangci-timeout                                | x       |          |
| golang-unit-tests-exclusions                    | x       |          |
| release-application-name                        |         |          |
| release-architecture                            |         |          |
| release-build-tags                              |         |          |
| release-dir                                     |         |          |
| release-type                                    |         |          |
| task-version                                    | x       |          |
| testing-type                                    |         |          |
| test-timeout                                    |         |          |
| token                                           |         |          |
| trivy-action-db                                 | x       |          |
| trivy-action-java-db                            | x       |          |

Note: If an **x** is registered in the Default column, refer to the
[action.yml](action.yml) for the corresponding value.

<!-- markdownlint-enable MD013 -->

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

Example:

```bash
# Fetch release information for the specified tag
RELEASE_RESPONSE=$(gh api \
  -H "Accept: application/vnd.github+json" \
  "/repos/$OWNER/$REPO/releases/tags/$TAG_NAME")

# Extract the release ID
RELEASE_ID=$(echo "$RELEASE_RESPONSE" | jq -r '.id')

echo "Release ID: $RELEASE_ID"
```

Step 3: Get the Asset ID

```bash
# Fetch the list of assets for the release
ASSETS_RESPONSE=$(gh api \
  -H "Accept: application/vnd.github+json" \
  "/repos/$OWNER/$REPO/releases/$RELEASE_ID/assets")

# Extract the asset ID for the specified asset name
ASSET_ID=$(echo "$ASSETS_RESPONSE" | jq -r \
  --arg NAME "$ASSET_NAME" \
  '.[] | select(.name == $NAME) | .id')

echo "Asset ID: $ASSET_ID"
```

Step 4: Download the Asset

```bash
# Download the asset using the asset ID
gh api \
  -H "Accept: application/octet-stream" \
  "/repos/$OWNER/$REPO/releases/assets/$ASSET_ID" > "$ASSET_NAME"

echo "Downloaded asset: $ASSET_NAME"
```
