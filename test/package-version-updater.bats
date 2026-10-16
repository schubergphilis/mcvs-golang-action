#!/usr/bin/env bats

# Requires mikefarah yq v4 and jq, as preinstalled on the GitHub ubuntu runner.
# Run with: bats test/

setup() {
  # shellcheck source=../scripts/package-version-updater.sh
  source "${BATS_TEST_DIRNAME}/../scripts/package-version-updater.sh"

  cd "${BATS_TEST_TMPDIR}"
  CALLS="${BATS_TEST_TMPDIR}/calls"
  export BUILD_TASKFILE="${BATS_TEST_TMPDIR}/task.yml"
  export DEPENDENCIES_LABEL=dependencies
  export PACKAGE_VERSION_UPDATER_BRANCH=package-version-updater

  cat >"${BUILD_TASKFILE}" <<'YAML'
vars:
  MOCKERY_VERSION: v3.1.0
  YQ_VERSION: v4.0.0
YAML

  cat >action.yml <<'YAML'
inputs:
  task-version:
    default: 3.1.0
YAML
}

# records the arguments of a stubbed command, one call per line
record() {
  echo "$*" >>"${CALLS}"
}

releases_json() {
  cat <<'JSON'
[
  {"tagName": "v2.10.0", "isDraft": false, "isPrerelease": false},
  {"tagName": "v1.9.9", "isDraft": false, "isPrerelease": false},
  {"tagName": "v2.9.0", "isDraft": false, "isPrerelease": false},
  {"tagName": "v2.11.0-rc1", "isDraft": false, "isPrerelease": false},
  {"tagName": "nightly", "isDraft": false, "isPrerelease": false},
  {"tagName": "v3.0.0-rc1", "isDraft": false, "isPrerelease": true},
  {"tagName": "v3.0.0", "isDraft": true, "isPrerelease": false}
]
JSON
}

@test "latest_stable_package_version_on_github picks the highest plain stable version" {
  gh() { releases_json; }

  run latest_stable_package_version_on_github example/repo

  [ "${status}" -eq 0 ]
  [ "${output}" = "v2.10.0" ]
}

@test "latest_stable_package_version_on_github honours the version prefix" {
  gh() { releases_json; }

  run latest_stable_package_version_on_github example/repo v1.

  [ "${status}" -eq 0 ]
  [ "${output}" = "v1.9.9" ]
}

@test "latest_stable_package_version_on_github prints null without a matching release" {
  gh() { releases_json; }

  run latest_stable_package_version_on_github example/repo v9.

  [ "${status}" -eq 0 ]
  [ "${output}" = "null" ]
}

@test "go_list_latest_version skips pre-releases" {
  go() { echo "golang.org/x/tools v0.9.0 v0.10.0 v0.11.0-rc1"; }

  run go_list_latest_version golang.org/x/tools

  [ "${output}" = "v0.10.0" ]
}

@test "latest_stable_package_versions exports found versions and skips failed lookups" {
  latest_stable_package_version_on_github() {
    case "$1" in
      mikefarah/yq) echo v4.1.0 ;;
      vektra/mockery) echo null ;;
      *) return 1 ;;
    esac
  }
  go_list_latest_version() { return 1; }

  latest_stable_package_versions 2>"${BATS_TEST_TMPDIR}/stderr"

  [ "${YQ_VERSION}" = "v4.1.0" ]
  [ -z "${MOCKERY_VERSION:-}" ]
  grep -q "failed to determine a version for mockery" "${BATS_TEST_TMPDIR}/stderr"
  [[ " ${FAILED_LOOKUPS[*]} " == *" mockery "* ]]
  [[ " ${FAILED_LOOKUPS[*]} " != *" yq "* ]]
}

@test "generate_pr_body_with_updates lists only changed versions" {
  export YQ_VERSION=v4.1.0
  export MOCKERY_VERSION=v3.1.0

  # set -u, as in main: versions of failed lookups are unset
  run bash -c "set -u; source '${BATS_TEST_DIRNAME}/../scripts/package-version-updater.sh'; generate_pr_body_with_updates >/dev/null; printf '%s' \"\${PR_BODY}\""

  [ "${status}" -eq 0 ]
  [ "${output}" = "Updated yq: v4.0.0 → v4.1.0" ]
}

@test "replace_versions_with_latest_stable_package_versions writes only found versions" {
  export YQ_VERSION=v4.1.0

  # set -u, as in main: versions of failed lookups are unset
  run bash -c "set -u; source '${BATS_TEST_DIRNAME}/../scripts/package-version-updater.sh'; replace_versions_with_latest_stable_package_versions"

  [ "${status}" -eq 0 ]

  [ "$(yq -r '.vars.YQ_VERSION' "${BUILD_TASKFILE}")" = "v4.1.0" ]
  [ "$(yq -r '.vars.MOCKERY_VERSION' "${BUILD_TASKFILE}")" = "v3.1.0" ]
}

@test "update_task_version updates action.yml and the PR body" {
  latest_stable_package_version_on_github() { record "$@"; echo v3.2.0; }
  PR_BODY=""

  update_task_version

  [ "$(cat "${CALLS}")" = "go-task/task v3." ]
  [ "$(yq -r '.inputs.task-version.default' action.yml)" = "3.2.0" ]
  [ "${PR_BODY}" = "Updated Task: 3.1.0 → 3.2.0"$'\n' ]
}

@test "update_task_version keeps action.yml when no release is found" {
  latest_stable_package_version_on_github() { echo null; }

  update_task_version 2>/dev/null

  [ "$(yq -r '.inputs.task-version.default' action.yml)" = "3.1.0" ]
  [ "${FAILED_LOOKUPS[*]}" = "Task" ]
}

@test "update_task_version fails when action.yml is not updated" {
  latest_stable_package_version_on_github() { echo v3.2.0; }
  # ignore in-place writes, so that the version check reads the old value
  yq() {
    if [[ "$*" == *" -i "* ]]; then
      return 0
    fi

    command yq "$@"
  }

  run update_task_version

  [ "${status}" -eq 1 ]
}

@test "checkout_branch_required_to_apply_package_version_updates rebuilds the branch from main" {
  git() { record git "$@"; }

  checkout_branch_required_to_apply_package_version_updates

  grep -qx "git checkout -B package-version-updater origin/main" "${CALLS}"
}

@test "github_labels creates or updates the label" {
  gh() { record gh "$@"; }

  github_labels

  grep -q "^gh label create dependencies .* --force$" "${CALLS}"
}

@test "commit_and_push_changes pushes the branch with a lease" {
  git() { record git "$@"; }

  commit_and_push_changes

  grep -qx "git push origin package-version-updater --force-with-lease" "${CALLS}"
}

@test "create_or_edit_pr edits an open PR" {
  gh() {
    record gh "$@"
    if [[ "$2" == "list" ]]; then echo 1; fi
  }
  PR_BODY="body"

  create_or_edit_pr

  grep -q "^gh pr edit package-version-updater " "${CALLS}"
  [ -z "$(grep "^gh pr create" "${CALLS}")" ]
}

@test "create_or_edit_pr creates a labelled PR when none is open" {
  gh() {
    record gh "$@"
    if [[ "$2" == "list" ]]; then echo 0; fi
  }
  PR_BODY="body"

  create_or_edit_pr

  grep -q "^gh pr create .*--label dependencies" "${CALLS}"
}

@test "main stops before committing when all versions are up to date" {
  latest_stable_package_versions() { :; }
  checkout_branch_required_to_apply_package_version_updates() { :; }
  generate_pr_body_with_updates() { :; }
  update_task_version() { :; }
  replace_versions_with_latest_stable_package_versions() { :; }
  git() { :; }
  github_labels() { record github_labels; }
  commit_and_push_changes() { record commit_and_push_changes; }
  create_or_edit_pr() { record create_or_edit_pr; }

  run main

  [ "${status}" -eq 0 ]
  [ ! -e "${CALLS}" ]
}

@test "main commits, pushes and opens a PR when versions changed" {
  latest_stable_package_versions() { :; }
  checkout_branch_required_to_apply_package_version_updates() { :; }
  generate_pr_body_with_updates() { :; }
  update_task_version() { :; }
  replace_versions_with_latest_stable_package_versions() { :; }
  git() { echo " M build/task.yml"; }
  github_labels() { record github_labels; }
  commit_and_push_changes() { record commit_and_push_changes; }
  create_or_edit_pr() { record create_or_edit_pr; }

  run main

  [ "${status}" -eq 0 ]
  [ "$(cat "${CALLS}")" = $'github_labels\ncommit_and_push_changes\ncreate_or_edit_pr' ]
}

@test "main pushes the found updates and then fails when a lookup failed" {
  latest_stable_package_versions() { FAILED_LOOKUPS+=(mockery); }
  checkout_branch_required_to_apply_package_version_updates() { :; }
  generate_pr_body_with_updates() { :; }
  update_task_version() { :; }
  replace_versions_with_latest_stable_package_versions() { :; }
  git() { echo " M build/task.yml"; }
  github_labels() { record github_labels; }
  commit_and_push_changes() { record commit_and_push_changes; }
  create_or_edit_pr() { record create_or_edit_pr; }

  run main

  [ "${status}" -eq 1 ]
  [[ "${output}" == *"failed to determine a version for: mockery"* ]]
  [ "$(cat "${CALLS}")" = $'github_labels\ncommit_and_push_changes\ncreate_or_edit_pr' ]
}
