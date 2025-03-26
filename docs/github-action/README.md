# GitHub action

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

## Configuration

Create a `.github/workflows/golang.yml` file with the following content:

```yml
---
name: Golang
"on": push
permissions:
  contents: read
  packages: read
jobs:
  mcvs-golang-action:
    strategy:
      matrix:
        args:
          - {
              release-architecture: "amd64",
              release-dir: "./cmd/path-to-app",
              release-type: "binary",
              release-application-name: "some-app",
            }
          - {
              release-architecture: "arm64",
              release-dir: "./cmd/path-to-app",
              release-type: "binary",
              release-application-name: "some-lambda-func",
              release-build-tags: "lambda.norpc",
            }
          - { testing-type: "component" }
          - { testing-type: "coverage" }
          - { testing-type: "integration" }
          - { testing-type: "lint", build-tags: "component" }
          - { testing-type: "lint", build-tags: "e2e" }
          - { testing-type: "lint", build-tags: "integration" }
          - { testing-type: "mcvs-texttidy" }
          - { testing-type: "security-golang-modules" }
          - { testing-type: "security-grype" }
          - { testing-type: "security-trivy", security-trivyignore: "" }
          - { testing-type: "unit" }
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
          test-timeout: ${{ env.test-timeout }}
          code-coverage-timeout: ${{ env.test-timeout }}
```

and a [.golangci.yml](https://golangci-lint.run/usage/configuration/).

| Option                                          | Default | Required |
| :---------------------------------------------- | :------ | -------- |
| build-tags                                      | x       |          |
| code-coverage-expected                          | x       |          |
| code-coverage-timeout                           |         |          |
| github-token-for-downloading-private-go-modules |         |          |
| golangci-timeout                                | x       |          |
| golang-unit-tests-exclusions                    | x       |          |
| grype-version                                   | x       |          |
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
