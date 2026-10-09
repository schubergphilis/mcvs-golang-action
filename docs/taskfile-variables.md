# Taskfile Variables

The `./build/task.yml` in this project defines a number of variables. Some of
these can be overridden when including this Taskfile in your project. See the
example below, where the `CODE_COVERAGE_STRICT` variable is overridden, for how
to do this.

The following variables can be overridden:

| Variable                    | Description                                                                                                          |
| :-------------------------- | :------------------------------------------------------------------------------------------------------------------- |
| `CODE_COVERAGE_STRICT`      | Enables or disables strict enforcement of setting the minimum coverage to the maximum observed coverage.             |
| `GOLANGCI_LINT_CONFIG_PATH` | Defines the path to the golangci-lint configuration file.                                                            |
| `GOTESTSUM_ENABLED`         | Enables or disables running the tests via [gotestsum](https://github.com/gotestyourself/gotestsum). Default: `true`. |
| `GOTESTSUM_FORMAT`          | The gotestsum [output format](https://github.com/gotestyourself/gotestsum#output-format). Default: `testname`.       |

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

Note: same goes for the `GOLANGCI_LINT_RUN_TIMEOUT_MINUTES` and
`GOLANGCI_LINT_INSTALL_ATTEMPTS` settings. The latter defaults to `3` and
bounds how many times the golangci-lint download is retried before the job
fails; set it to `1` to restore the previous fail-on-first-error behaviour.
