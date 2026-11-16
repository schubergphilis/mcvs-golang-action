# Caching

By default every job downloads all Go modules and compiles them from scratch.
On a project with many dependencies that is most of the duration of a job. Set
`go-cache: true` to restore the Go module cache, the Go build cache and the
golangci-lint cache.

## Usage

The caches are only saved by a run on the default branch, so the workflow has
to run there as well:

```yml
---
name: Golang
"on":
  pull_request:
  push:
    branches:
      - main
permissions:
  contents: read
  packages: read
jobs:
  MCVS-golang-action:
    strategy:
      matrix:
        args:
          - testing-type: "coverage"
          - testing-type: "lint"
          - testing-type: "security-golang-modules"
    runs-on: ubuntu-24.04
    steps:
      - uses: actions/checkout@v4.1.1
      - uses: schubergphilis/mcvs-golang-action@v3
        with:
          go-cache: true
          testing-type: ${{ matrix.args.testing-type }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

A pull request does not get faster until the default branch has run once with
`go-cache: true`.

## How it works

- **Restore**: the `component`, `coverage`, `integration`, `lint`,
  `mocks-tidy`, `security-golang-modules` and `unit` testing types restore the
  caches. The other testing types compile no Go code of the project.
- **Save**: only a run on the default branch saves, including a scheduled run.
  GitHub lets every pull request restore a cache of the default branch, but
  not one of another pull request. It also means that a pull request cannot
  write a cache that other runs restore.
- **Module cache**: one cache for all jobs, keyed on the `go.sum` files.
- **Build cache**: one cache per testing type and build tags, as those compile
  with different flags. It is keyed on the Go version and the `go.mod` and
  `go.sum` files. When those change, the most recent cache of the same job is
  restored and only what changed is compiled.
- **Tests are never served from the cache**: the build cache also holds test
  results. With `go-cache: true` the tests run with `-count=1`, so a test that
  depends on a service or a container image is not reported as passed without
  being run.
- **Release builds** (`release-dir` set) never restore a cache, so that a
  poisoned cache cannot end up in a published binary.

## Keeping the caches

GitHub removes a cache that has not been restored for seven days and limits
the caches of a repository to 10 GB. In a repository that is quiet for more
than a week, a scheduled run on the default branch keeps the caches:

```yml
"on":
  schedule:
    - cron: 17 4 * * 1,4
```
