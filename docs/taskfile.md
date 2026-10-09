# Taskfile

This repository offers configuration for [Task](https://taskfile.dev/): a
`./build/task.yml` which contains standard tasks, like installing and
running a linter.

This `./build/task.yml` can then be used by other projects. This has the
advantage that you do not need to copy and paste Makefile snippets from one
project to another. As a consequence each project using this `./build/task.yml`
immediately benefits from improvements made here (e.g. new tasks or
improvements in the tasks).

If you are new to Task, you may want to check out the following resources:

- [Installation instructions](https://taskfile.dev/installation/)
- Instructions to [configure completions](https://taskfile.dev/installation/#setup-completions)
- [Integrations](https://taskfile.dev/integrations/) with e.g. Visual Studio Code, Sublime and IntelliJ.
