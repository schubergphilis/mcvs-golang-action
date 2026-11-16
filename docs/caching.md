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
- **Save**: only a `push`, `schedule` or `workflow_dispatch` run on the
  default branch saves. GitHub lets every pull request restore a cache of the
  default branch, but not one of another pull request. It also means that a
  pull request cannot write a cache that other runs restore. Other events
  that run on the default branch, such as `pull_request_target`, never save,
  as a workflow may use those to test the code of a pull request.
- **Module cache**: one cache for all jobs, keyed on the `go.sum` files. It
  is saved by the first job that finishes, except `security-golang-modules`:
  that job downloads hardly any module source, and a cache that is saved once
  is not replaced until `go.sum` changes.
- **Build cache**: one cache per testing type and build tags, as those compile
  with different flags. It is keyed on the Go version and the `go.mod` and
  `go.sum` files. When those change, the most recent cache of the same job is
  restored and only what changed is compiled.
- **Tests are never served from the cache**: the build cache also holds test
  results. With `go-cache: true` the tests run with `-count=1`, so a test that
  depends on a service or a container image is not reported as passed without
  being run. It is added to the `GOFLAGS` that the runner or the container
  image already has.
- **Release builds** (`release-dir` set) never restore a cache, so that a
  poisoned cache cannot end up in a published binary.

## Private modules

Every pull request, also one from a fork, can restore a cache of the default
branch. The module cache holds the source of every module and the build cache
its compiled packages.

- When `github-token-for-downloading-private-go-modules` is set, nothing is
  cached, also with `go-cache: true`.
- Do not set `go-cache: true` in a public repository that downloads private
  modules in another way, for example with a `.netrc` file.

## Keeping the caches

GitHub removes a cache that has not been restored for seven days and limits
the caches of a repository to 10 GB. In a repository that is quiet for more
than a week, a scheduled run on the default branch keeps the caches:

```yml
"on":
  schedule:
    - cron: 17 4 * * 1,4
```
