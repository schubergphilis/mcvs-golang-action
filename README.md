# golang-action

Create a `.github/workflows/golang.yml` file with the following content:

```bash
---
name: Golang
'on': push
jobs:
  golang:
    runs-on: ubuntu-20.04
    steps:
      - uses: actions/checkout@v4.1.1
      - uses: schubergphilis/mcaf-mcvs-golang-action@v0.1.0
        with:
          inputs.golang-unit-tests-exclusions: "\(cmd\/some-app\|internal\/app\/some-app\)"
```
