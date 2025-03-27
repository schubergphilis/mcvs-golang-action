# MCVS-golang-action

The Mission Critical Vulnerability Scanner (MCVS) Golang Action repository
offers a set of standardized tools designed to maintain high-quality standards
in projects that use Golang code.

## Github action

Based on this
[GitHub Doc](https://docs.github.com/en/actions/sharing-automations/creating-actions/creating-a-composite-action),
a composite action has been created. Check
[this](./docs/github-action/README.md) to find out more.

## Taskfile

Another tool is configuration for [Task](https://taskfile.dev/). This
repository provides a `Taskfile.yml` that includes standard tasks, such as
installing and running a linter. Navigate to the
[docs](./docs/taskfile/README.md) for more details.

## Testing

### Integration

To execute integration tests, make sure that the code is located in a file with
a `_integration_test.go` postfix, such as `some_integration_test.go`.
Additionally, include the following header in the file:

```bash
//go:build integration
```

After adding this header, issue the command `task remote:test-integration --yes`
as demonstrated in this example. This action will run both unit and integration
tests. If `task remote:test --yes` is executed, only unit tests will be run.

### Component

See the integration paragraph for the steps and replace `integration` with
`component` to run them.

## Releases

### Downloading released assets from another private repository

You will need a personal access token (PAT) with the `repo` scope. To download
releases from a private repository. You can simply use the gh command or curl
to download the release assets. Please read the
[GitHub documentation](https://docs.github.com/en/rest/releases/assets)
for more information.
