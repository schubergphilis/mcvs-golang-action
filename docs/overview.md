# GitHub Action

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
- A test summary, including the number of tests that have been run per testing
  type, using [gotestsum](https://github.com/gotestyourself/gotestsum).

In summary, using this action will ensure that Golang code meets certain
standards before it will be deployed to production as the assembly line will
fail if an issue arises.

Note: there is an [internal action](../.github/workflows/package-version-updater.yml)
that will update package versions that cannot be updated by Dependabot.
