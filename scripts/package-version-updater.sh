#!/bin/bash

###############################################################################
# Weekly Dependency Version Updater Script
#
# PURPOSE:
#   Automates updating manually pinned package versions in a Taskfile (e.g., Taskfile.yml).
#   The script checks for the latest available versions, updates the Taskfile if needed,
#   and opens a pull request on GitHub summarizing all changes.
#   The PR body will only list packages whose versions have changed.
#
# REQUIREMENTS:
#   - Must set these environment variables:
#       BUILD_TASKFILE                Path to your Taskfile.yml
#       PACKAGE_VERSION_UPDATER_BRANCH  Name for the update branch
#       DEPENDENCIES_LABEL            Label name for dependency PRs (e.g. "dependencies")
#   - The following CLI tools must be available: gh, jq, yq, git, go, tr, grep
#
# USAGE:
#   Set the required environment variables and run this script in CI or locally:
#     export BUILD_TASKFILE=./Taskfile.yml
#     export PACKAGE_VERSION_UPDATER_BRANCH=package-version-updater
#     export DEPENDENCIES_LABEL=dependencies
#     ./update-script.sh
#
# CUSTOMIZATION: ADDING A NEW PACKAGE TO AUTO-UPDATE
# --------------------------------------------------
# To add or manage packages, simply update the PACKAGES_TO_BE_UPDATED array below.
#
#   Each entry in the array has the form:
#      "ENV_VAR YAML_VAR DISPLAY_NAME FETCH_FUNCTION FETCH_ARGUMENT"
#
#      ENV_VAR:         The name of the exported variable used in this script
#      YAML_VAR:        The variable name in the Taskfile's .vars section
#      DISPLAY_NAME:    Human-friendly name for PR changelog/report
#      FETCH_FUNCTION:  Name of the shell function to call for latest version
#      FETCH_ARGUMENT:  The argument passed to FETCH_FUNCTION (GitHub repo or module path)
#
# Example: To update a new package (e.g., example/examplepkg from GitHub),
#   1. Add a new entry to PACKAGES_TO_BE_UPDATED:
#        "EXAMPLEPKG_VERSION EXAMPLEPKG_VERSION examplepkg latest_stable_package_version_on_github example/examplepkg"
#
# That's it! No other code changes are required.
#
# TIP: Prefer short and readable display names.
###############################################################################

set -xeuo pipefail

readonly PACKAGES_TO_BE_UPDATED=(
  "GO_SWAGGER_VERSION GO_SWAGGER_VERSION go-swagger latest_stable_package_version_on_github go-swagger/go-swagger"
  "GQLGEN_VERSION GQLGEN_VERSION gqlgen latest_stable_package_version_on_github 99designs/gqlgen"
  "GQLGENC_VERSION GQLGENC_VERSION gqlgenc latest_stable_package_version_on_github Yamashou/gqlgenc"
  "GRAPHQL_LINTER_VERSION GRAPHQL_LINTER_VERSION graphql-linter latest_stable_package_version_on_github schubergphilis/graphql-linter"
  "MOCKERY_VERSION MOCKERY_VERSION mockery latest_stable_package_version_on_github vektra/mockery"
  "OPA_VERSION OPA_VERSION opa latest_stable_package_version_on_github open-policy-agent/opa"
  "PRESENT_VERSION PRESENT_VERSION present go_list_latest_version golang.org/x/tools"
  "REGAL_VERSION REGAL_VERSION regal latest_stable_package_version_on_github StyraOSS/regal"
  "YQ_VERSION YQ_VERSION yq latest_stable_package_version_on_github mikefarah/yq"
)
readonly PR_TITLE="build(deps): weekly update package versions that cannot be updated by dependabot"

check_label_exists() {
  local label_name="$1"

  LABEL_EXISTS=$(
    gh label list --json name |
    jq -r '
      .[] |
      select(.name == "'"$label_name"'") |
      .name
    '
  )
  if [ -z "${LABEL_EXISTS}" ]; then
    echo "label: '${label_name}' does NOT exist"
    return 1
  fi
}

latest_stable_package_version_on_github() {
  gh release list \
    --repo $1 \
    --limit 100 \
    --json tagName,isDraft,isPrerelease,publishedAt | \
      jq -r '[.[] | select(.isDraft == false and .isPrerelease == false)] | sort_by(.publishedAt) | last.tagName'
}

latest_stable_package_versions() {
  for dep in "${PACKAGES_TO_BE_UPDATED[@]}"; do
    set -- $dep
    local env_var="$1"
    local yaml_var="$2"
    local display_name="$3"
    local fetch_func="$4"
    local fetch_arg="$5"

    local version
    version="$($fetch_func "$fetch_arg")"
    export "$env_var"="$version"
    echo "$env_var: $version"
  done
}

go_list_latest_version() {
  local module="$1"
  go list -m -versions "$module" | \
    tr ' ' '\n' | \
    grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | \
    tail -1
}

checkout_branch_required_to_apply_package_version_updates() {
  git fetch -p -P

  if (git ls-remote --exit-code --heads origin refs/heads/${PACKAGE_VERSION_UPDATER_BRANCH}); then
    echo "Branch '${PACKAGE_VERSION_UPDATER_BRANCH}' already exists."
    git checkout ${PACKAGE_VERSION_UPDATER_BRANCH}

    return
  fi

  git checkout -b ${PACKAGE_VERSION_UPDATER_BRANCH}
}

replace_versions_with_latest_stable_package_versions() {
  for dep in "${PACKAGES_TO_BE_UPDATED[@]}"; do
    set -- $dep

    local env_var="$1"
    local yaml_var="$2"
    local display_name="$3"
    local version="${!env_var}"

    echo "$env_var: $version"
    yq -i ".vars.${yaml_var} = strenv(${env_var})" "$BUILD_TASKFILE"
  done
}

github_labels() {
  if ! check_label_exists ${DEPENDENCIES_LABEL}; then
    gh label create "${DEPENDENCIES_LABEL}" \
      --color "#0366d6" \
      --description "Pull requests that update a dependency file"
  fi

  labels=("${DEPENDENCIES_LABEL}")
  echo "Labels:"

  for label in "${labels[@]}"; do
    echo "'$label'"
  done
}

commit_and_push_changes() {
  if [ -n "$(git status --porcelain)" ]; then echo "There are uncommitted changes."; else echo "No changes to commit." && return; fi
    git add ${BUILD_TASKFILE}
    git config user.name github-actions[bot]
    git config user.email 41898282+github-actions[bot]@users.noreply.github.com

  if ! git commit -m "${PR_TITLE}"; then git commit --amend --no-edit; fi
    git push origin ${PACKAGE_VERSION_UPDATER_BRANCH} --force-with-lease
}

create_or_edit_pr() {
  if gh pr list --json title | jq -e '.[] | select(.title | test("build\\(deps\\): weekly update package versions that cannot be updated by dependabot"))'; then
    echo "PR exists already. Updating the 'title' and 'description'..."

    gh pr edit ${PACKAGE_VERSION_UPDATER_BRANCH} \
      --body "${PR_BODY}" \
      --title "${PR_TITLE}"

    return
  fi

  echo "creating pr..."
  label_args=()
  for label in "${labels[@]}"; do
    label_args+=(--label "$label")
  done

  gh pr create \
    --base main \
    --body "${PR_BODY}" \
    --fill \
    --head "${PACKAGE_VERSION_UPDATER_BRANCH}" \
    --title "${PR_TITLE}" \
    "${label_args[@]}"
}

generate_pr_body_with_updates() {
  local pr_body=""

  for dep in "${PACKAGES_TO_BE_UPDATED[@]}"; do
    set -- $dep

    local env_var="$1"
    local yaml_var="$2"
    local display_name="$3"
    local new_version="${!env_var}"
    local old_version

    old_version=$(yq -r ".vars.${yaml_var}" "$BUILD_TASKFILE")

    if [[ -z "$new_version" || -z "$old_version" ]]; then
      continue
    fi

    if [[ "$new_version" != "$old_version" ]]; then
      pr_body+="Updated $display_name: $old_version → $new_version"$'\n'
    fi
  done

  if [[ -z "$pr_body" ]]; then
    pr_body="No dependency versions were updated."
  fi

  export PR_BODY="$pr_body"
  echo "PR_BODY: ${PR_BODY}"
}

main() {
  latest_stable_package_versions
  checkout_branch_required_to_apply_package_version_updates
  generate_pr_body_with_updates
  replace_versions_with_latest_stable_package_versions
  github_labels
  commit_and_push_changes
  create_or_edit_pr
}

main
