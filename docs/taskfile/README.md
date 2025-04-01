# Taskfile

This `Taskfile.yml` can then be used by other projects. This has the advantage
that you do not need to copy and paste Makefile snippets from one project to
another. Consequently, each project using this `Taskfile.yml` immediately
benefits from improvements made here, such as new tasks or enhancements in the
existing tasks.

If you are new to Task, you may want to check out the following resources:

- [Installation instructions](https://taskfile.dev/installation/)
- Instructions to [configure completions](https://taskfile.dev/installation/#setup-completions)
- [Integrations](https://taskfile.dev/integrations/) with e.g. Visual Studio
  Code, Sublime and IntelliJ.

### Configuration

The `Taskfile.yml` in this project defines a number of variables. Some of these
can be overridden when including this Taskfile in your project. See the example
below, where the `MOCKERY_VERSION` variable is overridden, for how to do this.

The following variables can be overridden:

| Variable          | Description                |
| :---------------- | :------------------------- |
| `MOCKERY_VERSION` | Define the Mockery version |

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
      MOCKERY_VERSION: v1.2.3
```

Note: same goes for the `GOLANGCI_LINT_RUN_TIMEOUT_MINUTES` setting.
