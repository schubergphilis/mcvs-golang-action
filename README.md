# MCVS-golang-action

Mission Critical Vulnerability Scanner (MCVS) Golang Action is a custom
[GitHub Action](https://github.com/features/actions) that consists of the
following steps:

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

## Usage

Create a `.github/workflows/golang.yml` file with the following content:

```yaml
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
        testing-type:
          - component
          - coverage
          - integration
          - lint
          - security-golang-modules
          - security-grype
          - security-trivy
          - unit
    runs-on: ubuntu-22.04
    steps:
      - uses: actions/checkout@v4.1.1
      - uses: schubergphilis/mcvs-golang-action@v0.9.0
        with:
          golang-unit-tests-exclusions: |-
            \(cmd\/some-app\|internal\/app\/some-app\)
          testing-type: ${{ matrix.testing-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

and a [.golangci.yml](https://golangci-lint.run/usage/configuration/).

<!-- markdownlint-disable MD013 -->

| Option                             | Default                              | Required | Description                                                                                                      |
| :--------------------------------- | :----------------------------------- | -------- | :--------------------------------------------------------------------------------------------------------------- |
| code_coverage_expected             | 80                                   |          |                                                                                                                  |
| gci                                | true                                 |          | Check for 'incorrect import order'. If failed then instructions are shown to resolve the issue                   |
| golang-unit-tests-exclusions       | ' '                                  |          |                                                                                                                  |
| golangci-lint-version              | v1.55.2                              |          |                                                                                                                  |
| golang-number-of-tests-in-parallel | 4                                    |          |                                                                                                                  |
| token                              | ' '                                  | x        | GitHub token that is required to push an image to the registry of the project and to pull cached Trivy DB images |
| trivy-action-db                    | ghcr.io/aquasecurity/trivy-db:2      |          | Replace this with a cached image to prevent bump into pull rate limiting issues                                  |
| trivy-action-java-db               | ghcr.io/aquasecurity/trivy-java-db:1 |          | Replace this with a cached image to prevent bump into pull rate limiting issues                                  |

<!-- markdownlint-enable MD013 -->

### Integration

To execute integration tests, make sure that the code is located in a file with
a `_integration_test.go` postfix, such as `some_integration_test.go`.
Additionally, include the following header in the file:

```bash
//go:build integration
```

After adding this header, issue the command `go test ./... --tags=integration`
as demonstrated in this example. This action will run both unit and integration
tests. If the `--tags` step is omitted, only unit tests will be executed.

### Component

See the integration paragraph for the steps and replace `integration` with
`component` to run them.
