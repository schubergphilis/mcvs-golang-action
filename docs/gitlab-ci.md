# GitLab CI

The checks of the GitHub Action are also available for GitLab CI as two
[CI/CD components](https://docs.gitlab.com/ci/components/):

- [templates/golang.yml](../templates/golang.yml): one job per include that
  runs the checks of the given `testing-type`, like the action does.
- [templates/golang-release.yml](../templates/golang-release.yml): builds a
  binary, like the `release-*` inputs of the action do.

Both use the remote Taskfile, so the project needs the `Taskfile.yml` that is
described in the [Quickstart](../README.md#locally).

## Including the components

This repository lives on GitHub, so it is not in a GitLab CI/CD Catalog. Until
it is mirrored to a GitLab project, include the templates with `include:remote`,
which accepts the same `inputs` as a component:

```yml
---
include:
  - remote: https://raw.githubusercontent.com/schubergphilis/mcvs-golang-action/<ref>/templates/golang.yml
    inputs:
      testing-type: unit
  - remote: https://raw.githubusercontent.com/schubergphilis/mcvs-golang-action/<ref>/templates/golang.yml
    inputs:
      testing-type: lint
  - remote: https://raw.githubusercontent.com/schubergphilis/mcvs-golang-action/<ref>/templates/golang.yml
    inputs:
      build-tags: integration
      job-name-suffix: -integration
      testing-type: lint
  - remote: https://raw.githubusercontent.com/schubergphilis/mcvs-golang-action/<ref>/templates/golang.yml
    inputs:
      code-coverage-expected: "80"
      testing-type: coverage
  - remote: https://raw.githubusercontent.com/schubergphilis/mcvs-golang-action/<ref>/templates/golang.yml
    inputs:
      testing-type: security-golang-modules
```

Use a tag or commit SHA as `<ref>`. When the repository is mirrored to a GitLab
project, `include:component` works with the same inputs:

```yml
include:
  - component: $CI_SERVER_FQDN/<group>/mcvs-golang-action/golang@<version>
    inputs:
      testing-type: unit
```

Every include adds a job named `mcvs-golang-<testing-type>`, followed by the
`job-name-suffix` input. Use the suffix to include the same testing type twice,
like the second `lint` job above.

## Inputs

`templates/golang.yml` accepts the inputs of the action with the same names
and defaults: `build-tags`, `code-coverage-expected`,
`code-coverage-opa-expected`, `code-coverage-timeout`,
`golang-unit-tests-exclusions`, `golangci-timeout`, `gta-base-branch`,
`gta-partial-testing`, `task-version`, `test-timeout` and `testing-type`. The
`testing-type` values are the same as for the action. In addition:

- `image`: the image of the job, `golang:1` by default. `GOTOOLCHAIN=auto`
  makes go download the version that is defined in `go.mod`.
- `job-name-suffix`: see above.
- `stage`: the stage of the job, `test` by default.

`security-grype` runs [Grype](https://github.com/anchore/grype) through
`task remote:grype`, with the version that is pinned in the Taskfile.

Tokens are not inputs, because inputs end up in the pipeline configuration.
Set them as masked CI/CD variables instead:

- `GITHUB_TOKEN_FOR_DOWNLOADING_PRIVATE_GO_MODULES`: a token that can read the
  private Go modules on github.com. It is set through `GIT_CONFIG_*`
  environment variables, so it is never written to disk.
- `GITHUB_TOKEN`: used by some tools to avoid GitHub API rate limits.

## Releases

`templates/golang-release.yml` accepts the `release-application-name`,
`release-architecture`, `release-build-tags`, `release-dir`, `release-os` and
`release-type` inputs of the action, plus `image`, `job-name-suffix` and
`stage` (`build` by default):

```yml
include:
  - remote: https://raw.githubusercontent.com/schubergphilis/mcvs-golang-action/<ref>/templates/golang-release.yml
    inputs:
      job-name-suffix: -amd64
      release-application-name: mcvs-image-downloader
      release-architecture: amd64
      release-dir: cmd/mcvs-image-downloader
      release-type: binary
```

In a tag pipeline, the binary is uploaded to the generic package registry,
which has to be enabled for the project, and linked to the release of that
tag. The release is created if it does not exist yet. Both use `CI_JOB_TOKEN`.
