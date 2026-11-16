#!/usr/bin/env bats

# Run with: bats test/

setup() {
  # shellcheck source=../scripts/go-cache-config.sh
  source "${BATS_TEST_DIRNAME}/../scripts/go-cache-config.sh"
}

# stands in for the go CLI, so the tests do not depend on the installed version
go() {
  case "$2" in
    GOCACHE) echo /home/runner/.cache/go-build ;;
    GOMODCACHE) echo /home/runner/go/pkg/mod ;;
    GOVERSION) echo go1.27.4 ;;
  esac
}

@test "go_cache_save: a push to the default branch saves" {
  run go_cache_save push refs/heads/main main
  [ "${output}" = "true" ]
}

@test "go_cache_save: a workflow_dispatch on the default branch saves" {
  run go_cache_save workflow_dispatch refs/heads/main main
  [ "${output}" = "true" ]
}

@test "go_cache_save: a scheduled run saves, also without a default branch" {
  run go_cache_save schedule refs/heads/main ""
  [ "${output}" = "true" ]
}

@test "go_cache_save: a default branch with another name saves" {
  run go_cache_save push refs/heads/master master
  [ "${output}" = "true" ]
}

@test "go_cache_save: a pull request does not save" {
  run go_cache_save pull_request refs/pull/1/merge main
  [ "${output}" = "false" ]
}

@test "go_cache_save: a push to another branch does not save" {
  run go_cache_save push refs/heads/feature main
  [ "${output}" = "false" ]
}

@test "go_cache_save: a branch that only starts with the default branch does not save" {
  run go_cache_save push refs/heads/main-2 main
  [ "${output}" = "false" ]
}

@test "go_cache_save: a push of a tag does not save" {
  run go_cache_save push refs/tags/v1.0.0 main
  [ "${output}" = "false" ]
}

@test "go_cache_save: a workflow_dispatch on another branch does not save" {
  run go_cache_save workflow_dispatch refs/heads/feature main
  [ "${output}" = "false" ]
}

@test "go_cache_save: a push without a known default branch does not save" {
  run go_cache_save push refs/heads/ ""
  [ "${output}" = "false" ]
}

# These events run on the default branch, but a workflow may use them to test
# the code of a pull request.
@test "go_cache_save: pull_request_target does not save" {
  run go_cache_save pull_request_target refs/heads/main main
  [ "${output}" = "false" ]
}

@test "go_cache_save: issue_comment does not save" {
  run go_cache_save issue_comment refs/heads/main main
  [ "${output}" = "false" ]
}

@test "go_cache_save: workflow_run does not save" {
  run go_cache_save workflow_run refs/heads/main main
  [ "${output}" = "false" ]
}

@test "go_cache_save: merge_group does not save" {
  run go_cache_save merge_group refs/heads/gh-readonly-queue/main/pr-1 main
  [ "${output}" = "false" ]
}

@test "go_cache_save_modules: a job that runs go mod tidy saves the modules" {
  run go_cache_save_modules true coverage
  [ "${output}" = "true" ]
}

@test "go_cache_save_modules: security-golang-modules does not save the modules" {
  run go_cache_save_modules true security-golang-modules
  [ "${output}" = "false" ]
}

@test "go_cache_save_modules: nothing saves the modules when the run may not save" {
  run go_cache_save_modules false coverage
  [ "${output}" = "false" ]
}

@test "go_cache_variant: joins the testing type and the build tags" {
  run go_cache_variant lint component
  [ "${output}" = "lint-component" ]
}

@test "go_cache_variant: a testing type without build tags ends with a dash" {
  run go_cache_variant unit ""
  [ "${output}" = "unit-" ]
}

@test "go_cache_variant: replaces what a cache key or an output cannot hold" {
  run go_cache_variant lint $'component,e2e integration\ninjected=1'
  [ "${output}" = "lint-component-e2e-integration-injected-1" ]
}

@test "main: writes the outputs for a pull request" {
  export BUILD_TAGS=integration DEFAULT_BRANCH=main TESTING_TYPE=lint
  export GITHUB_EVENT_NAME=pull_request GITHUB_REF=refs/pull/1/merge
  export GITHUB_OUTPUT="${BATS_TEST_TMPDIR}/output" RUNNER_ARCH=X64 RUNNER_OS=Linux

  main

  run cat "${GITHUB_OUTPUT}"
  [ "${lines[0]}" = "build-cache=/home/runner/.cache/go-build" ]
  [ "${lines[1]}" = "build-key=go-build-Linux-X64-go1.27.4-lint-integration-" ]
  [ "${lines[2]}" = "enabled=true" ]
  [ "${lines[3]}" = "module-cache=/home/runner/go/pkg/mod" ]
  [ "${lines[4]}" = "module-key=go-modules-Linux-X64-" ]
  [ "${lines[5]}" = "save=false" ]
  [ "${lines[6]}" = "save-modules=false" ]
}

@test "main: a push to the default branch saves both caches" {
  export BUILD_TAGS="" DEFAULT_BRANCH=main TESTING_TYPE=coverage
  export GITHUB_EVENT_NAME=push GITHUB_REF=refs/heads/main
  export GITHUB_OUTPUT="${BATS_TEST_TMPDIR}/output" RUNNER_ARCH=X64 RUNNER_OS=Linux

  main

  run cat "${GITHUB_OUTPUT}"
  [ "${lines[5]}" = "save=true" ]
  [ "${lines[6]}" = "save-modules=true" ]
}
