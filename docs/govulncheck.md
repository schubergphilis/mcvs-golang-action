# govulncheck

Add a `.govulncheck.yaml` file to a project to exclude certain vulnerabilities
that cannot be fixed right away. Currently, it is
[not supported](https://github.com/golang/go/issues/59507).

```zsh
---
vulnerability:
  exclude:
    - GO-2025-4020 # https://github.com/anchore/syft/issues/4338
```
