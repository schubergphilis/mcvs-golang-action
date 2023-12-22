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
| golang-number-of-tests-in-parallel | 1       |
