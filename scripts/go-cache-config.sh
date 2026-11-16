#!/bin/bash

###############################################################################
# Go cache configuration
#
# PURPOSE:
#   Decides for the 'go-cache' input of action.yml which caches a job restores
#   and whether it may save them, and writes that to GITHUB_OUTPUT.
#
#   Whether a job may save is the security gate of the feature: every pull
#   request restores the caches of the default branch, so the code of a pull
#   request must never end up in them.
#
# REQUIREMENTS:
#   - These environment variables, which GitHub Actions sets:
#       GITHUB_EVENT_NAME, GITHUB_OUTPUT, GITHUB_REF, RUNNER_ARCH, RUNNER_OS
#   - These environment variables, which action.yml sets:
#       BUILD_TAGS      The build-tags input
#       DEFAULT_BRANCH  The default branch of the repository
#       TESTING_TYPE    The testing-type input
#   - The go CLI.
#
# USAGE:
#   bash scripts/go-cache-config.sh
###############################################################################

# Jobs that compile differently get a build cache of their own. A cache key
# cannot contain a comma, which build tags are separated by.
go_cache_variant() {
  local testing_type="$1"
  local build_tags="$2"

  printf '%s' "${testing_type}-${build_tags}" | tr -c 'A-Za-z0-9._-' '-'
}

# Only a run on the default branch saves: every pull request can restore its
# caches, and a pull request cannot write a cache that other runs restore. A
# scheduled run always runs on the default branch.
#
# The event is checked as well as the ref. pull_request_target, issue_comment
# and workflow_run also run on the default branch, and a workflow may use
# those to test the code of a pull request.
go_cache_save() {
  local event_name="$1"
  local ref="$2"
  local default_branch="$3"

  case "${event_name}" in
    schedule)
      echo true
      ;;
    push | workflow_dispatch)
      if [[ -n "${default_branch}" && "${ref}" == "refs/heads/${default_branch}" ]]; then
        echo true
      else
        echo false
      fi
      ;;
    *)
      echo false
      ;;
  esac
}

# All jobs share one module cache and the first job that finishes saves it.
# security-golang-modules does not run 'go mod tidy' and 'go mod verify'
# downloads no module sources, so its module cache is nearly empty. If it
# saved, the other jobs would get an exact hit and never save a complete one.
go_cache_save_modules() {
  local save="$1"
  local testing_type="$2"

  if [[ "${save}" == "true" && "${testing_type}" != "security-golang-modules" ]]; then
    echo true
  else
    echo false
  fi
}

main() {
  set -euo pipefail

  local go_version save variant

  variant=$(go_cache_variant "${TESTING_TYPE:-}" "${BUILD_TAGS:-}")
  # A build cache is of no use to another Go version.
  go_version=$(go env GOVERSION)
  save=$(go_cache_save "${GITHUB_EVENT_NAME:-}" "${GITHUB_REF:-}" "${DEFAULT_BRANCH:-}")

  {
    echo "build-cache=$(go env GOCACHE)"
    echo "build-key=go-build-${RUNNER_OS}-${RUNNER_ARCH}-${go_version}-${variant}-"
    echo "enabled=true"
    echo "module-cache=$(go env GOMODCACHE)"
    echo "module-key=go-modules-${RUNNER_OS}-${RUNNER_ARCH}-"
    echo "save=${save}"
    echo "save-modules=$(go_cache_save_modules "${save}" "${TESTING_TYPE:-}")"
  } >>"${GITHUB_OUTPUT}"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  main "$@"
fi
