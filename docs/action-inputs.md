# Action Inputs

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
          - testing-type: "unit"
    runs-on: ubuntu-24.04
    env:
      test-timeout: 10m0s
    steps:
      - uses: actions/checkout@v4.1.1
        with:
          fetch-depth: 0 # this is necessary for gta partial testing
      - uses: schubergphilis/mcvs-golang-action@v3
        with:
          build-tags: ${{ matrix.args.build-tags }}
          golang-unit-tests-exclusions: |-
            \(cmd\/some-app\|internal\/app\/some-app\)
          gta-base-branch: main
          gta-partial-testing: true
          release-architecture: ${{ matrix.args.release-architecture }}
          release-dir: ${{ matrix.args.release-dir }}
          release-type: ${{ matrix.args.release-type }}
          task-install: yes
          testing-type: ${{ matrix.args.testing-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
          test-timeout: ${{ env.test-timeout }}
          code-coverage-timeout: ${{ env.test-timeout }}
```

and a [.golangci.yml](https://golangci-lint.run/usage/configuration/).

<!-- markdownlint-disable MD013 -->

| Option                                          | Default | Required | Description                                                                                         |
| :---------------------------------------------- | :------ | -------- | :-------------------------------------------------------------------------------------------------- |
| build-tags                                      |         |          | Build tags to use when running tests and linting (e.g., "integration", "component", "e2e")          |
| code-coverage-expected                          | x       |          | Minimum expected code coverage percentage for standard tests                                        |
| code-coverage-opa-expected                      | x       |          | Minimum expected code coverage percentage for OPA (Open Policy Agent) tests                         |
| code-coverage-timeout                           |         |          | Timeout duration for code coverage analysis (e.g., "10m0s")                                         |
| github-token-for-downloading-private-go-modules |         |          | GitHub token with permissions to download Go modules from private repositories                      |
| go-version-file                                 | x       |          | Path to the go.mod or go.work file used to determine the Go version                                 |
| golangci-timeout                                |         |          | Timeout duration for golangci-lint execution                                                        |
| golang-unit-tests-exclusions                    | x       |          | Regex pattern to exclude specific packages from unit testing (e.g., `\(cmd\/app\|internal\/app\)`)  |
| grype-version                                   |         |          | Specific version of Grype vulnerability scanner to use                                              |
| gta-base-branch                                 | x       |          | The branch changed go packages will be compared to, to perform partial tests                        |
| gta-partial-testing                             | x       |          | Whether to run partial tests (true or false)                                                        |
| release-application-name                        |         |          | Name of the application binary to build (required when release-type is set)                         |
| release-architecture                            |         |          | Target architecture for the binary (e.g., "amd64", "arm64")                                         |
| release-build-tags                              |         |          | Build tags to use when building the release binary (e.g., "lambda.norpc")                           |
| release-dir                                     |         |          | Directory containing the main.go file for the binary to build                                       |
| release-os                                      | x       |          | Target operating system for the binary (e.g., "linux", "darwin")                                    |
| release-type                                    |         |          | Type of release to build (e.g., "binary")                                                           |
| task-install                                    | x       |          | Whether to install Task runner ("yes" or "no")                                                      |
| task-version                                    | x       |          | Version of Task runner to install                                                                   |
| testing-type                                    |         |          | Type of testing to run (e.g., "unit", "integration", "lint", "coverage", "security-golang-modules") |
| test-timeout                                    |         |          | Timeout duration for test execution (e.g., "10m0s")                                                 |
| token                                           |         |          | GitHub token for authentication (typically ${{ secrets.GITHUB_TOKEN }})                             |

Note: If an **x** is registered in the Default column, refer to the
[action.yml](../action.yml) for the corresponding value.

<!-- markdownlint-enable MD013 -->
