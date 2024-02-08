# golang-action

Create a `.github/workflows/golang.yml` file with the following content:

```bash
---
name: Golang
'on': push
jobs:
  mvcs-golang-action:
    runs-on: ubuntu-20.04
    steps:
      - uses: actions/checkout@v4.1.1
      - uses: schubergphilis/mcvs-golang-action@v0.1.1
        with:
          golang-unit-tests-exclusions: |-
            \(cmd\/some-app\|internal\/app\/some-app\)
```

and a [.golangci.yml](https://golangci-lint.run/usage/configuration/).

| option                             | default |
| ---------------------------------- | ------- |
| golang-unit-tests-exclusions       | ' '     |
| golangci-lint-version              | v1.55.2 |
| golang-number-of-tests-in-parallel | 4       |

## integration

In order to run integration tests ensure that the code resides in a file with
a `_integration_test.go` postfix, e.g., `some_integration_test.go` and add the
following header:

```bash
//go:build integration
```

Once this has been added and `go test ./... --tags=integration` is issued like
in this action, both unit and integration tests will be run. If the `--tags`
step is omitted, only unit tests will be run.
