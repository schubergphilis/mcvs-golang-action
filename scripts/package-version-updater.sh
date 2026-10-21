#!/bin/bash

###############################################################################
# Weekly Dependency Version Updater Script
#
# PURPOSE:
#   Automates updating manually pinned package versions in a Taskfile (e.g., Taskfile.yml),
#   the Task version in action.yml and templates/golang.yml, and the golang image of
#   the GitLab CI/CD components in templates/. The script checks for the latest
#   available versions, rebuilds the update branch from main, applies them and opens
#   or updates a pull request on GitHub summarizing all changes. Nothing is pushed if
#   all versions are up to date.
#   The PR body will only list packages whose versions have changed.
#
# REQUIREMENTS:
#   - Must set these environment variables:
#       BUILD_TASKFILE                  Path to your Taskfile.yml
#       PACKAGE_VERSION_UPDATER_BRANCH  Name for the update branch
#       DEPENDENCIES_LABEL              Label name for dependency PRs (e.g. "dependencies")
#   - The following CLI tools must be available: gh, jq, yq, git, go, tr, grep, curl, shasum
#
# USAGE:
#   Set the required environment variables and run this script in CI or locally:
#     export BUILD_TASKFILE=./Taskfile.yml
#     export PACKAGE_VERSION_UPDATER_BRANCH=package-version-updater
#     export DEPENDENCIES_LABEL=dependencies
#     ./scripts/package-version-updater.sh
#
# CUSTOMIZATION: ADDING A NEW PACKAGE TO AUTO-UPDATE
# --------------------------------------------------
# To add or manage packages, simply update the PACKAGES_TO_BE_UPDATED array below.
#
#   Each entry in the array has the form:
#      "ENV_VAR YAML_VAR DISPLAY_NAME FETCH_FUNCTION FETCH_ARGUMENT [VERSION_PREFIX]"
#
#      ENV_VAR:         The name of the exported variable used in this script
#      YAML_VAR:        The variable name in the Taskfile's .vars section
#      DISPLAY_NAME:    Human-friendly name for PR changelog/report
#      FETCH_FUNCTION:  Name of the shell function to call for latest version
#      FETCH_ARGUMENT:  The argument passed to FETCH_FUNCTION (GitHub repo or module path)
#      VERSION_PREFIX:  Optional tag prefix filter (for example: v3.)
#
# Example: To update a new package (e.g., example/examplepkg from GitHub),
#   1. Add a new entry to PACKAGES_TO_BE_UPDATED:
#        "EXAMPLEPKG_VERSION EXAMPLEPKG_VERSION examplepkg latest_stable_package_version_on_github example/examplepkg"
#
# That's it! No other code changes are required.
#
# TIP: Prefer short and readable display names.
###############################################################################

readonly PACKAGES_TO_BE_UPDATED=(
  "GOLANGCI_LINT_VERSION GOLANGCI_LINT_VERSION golangci-lint latest_stable_package_version_on_github golangci/golangci-lint"
  "GOLINES_VERSION GOLINES_VERSION golines latest_stable_package_version_on_github segmentio/golines"
  "GO_SWAGGER_VERSION GO_SWAGGER_VERSION go-swagger latest_stable_package_version_on_github go-swagger/go-swagger"
  "GQLGEN_VERSION GQLGEN_VERSION gqlgen latest_stable_package_version_on_github 99designs/gqlgen"
  "GQLGENC_VERSION GQLGENC_VERSION gqlgenc latest_stable_package_version_on_github Yamashou/gqlgenc"
  "GRAPHQL_LINTER_VERSION GRAPHQL_LINTER_VERSION graphql-linter latest_stable_package_version_on_github schubergphilis/graphql-linter"
  "GRYPE_VERSION GRYPE_VERSION grype latest_stable_package_version_on_github anchore/grype"
  "GOTESTSUM_VERSION GOTESTSUM_VERSION gotestsum latest_stable_package_version_on_github gotestyourself/gotestsum"
  "MOCKERY_VERSION MOCKERY_VERSION mockery latest_stable_package_version_on_github vektra/mockery v3."
  "OPA_VERSION OPA_VERSION opa latest_stable_package_version_on_github open-policy-agent/opa"
  "OSV_SCANNER_VERSION OSV_SCANNER_VERSION osv-scanner latest_stable_package_version_on_github google/osv-scanner"
  "PRESENT_VERSION PRESENT_VERSION present go_list_latest_version golang.org/x/tools"
  "REGAL_VERSION REGAL_VERSION regal latest_stable_package_version_on_github StyraOSS/regal"
  "YQ_VERSION YQ_VERSION yq latest_stable_package_version_on_github mikefarah/yq"
  "GTA_VERSION GTA_VERSION gta latest_stable_package_version_on_github digitalocean/gta"
  "WSL_VERSION WSL_VERSION wsl latest_stable_package_version_on_github bombsimon/wsl"
)
# packages whose version lookup failed; main fails the run at the end if any
FAILED_LOOKUPS=()
readonly PR_TITLE="build(deps): weekly update package versions that cannot be updated by dependabot"

latest_stable_package_version_on_github() {
  local repository="$1"
  local version_prefix="${2:-}"

  # only plain vX.Y.Z tags, so an rc tag that is not marked as a prerelease is
  # skipped, sorted by version rather than publish date, so a backport release
  # does not cause a downgrade
  gh release list \
    --repo "${repository}" \
    --limit 100 \
    --json tagName,isDraft,isPrerelease | \
      jq -r --arg version_prefix "${version_prefix}" '[.[] | select(.isDraft == false and .isPrerelease == false and (.tagName | startswith($version_prefix)) and (.tagName | test("^v?[0-9]+(\\.[0-9]+)*$")))] | sort_by(.tagName | ltrimstr("v") | split(".") | map(tonumber)) | last.tagName'
}

latest_stable_package_versions() {
  local dep

  for dep in "${PACKAGES_TO_BE_UPDATED[@]}"; do
    set -- $dep
    local env_var="$1"
    local yaml_var="$2"
    local display_name="$3"
    local fetch_func="$4"
    local fetch_arg="$5"
    local version_prefix="${6:-}"

    local version
    version="$($fetch_func "$fetch_arg" "$version_prefix")" || version=""

    if [[ -z "${version}" || "${version}" == "null" ]]; then
      echo "Warning: failed to determine a version for ${display_name}" >&2
      FAILED_LOOKUPS+=("${display_name}")
      continue
    fi

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

  # always rebuild the branch from main, so that the old versions in the PR
  # body are the ones on main and the branch picks up new commits from main.
  # This drops manual commits on the branch, and a failed lookup drops that
  # package's earlier bump from an open PR until a run succeeds; the run then
  # fails, so it is visible.
  git checkout -B "${PACKAGE_VERSION_UPDATER_BRANCH}" origin/main
}

replace_versions_with_latest_stable_package_versions() {
  local dep

  for dep in "${PACKAGES_TO_BE_UPDATED[@]}"; do
    set -- $dep

    local env_var="$1"
    local yaml_var="$2"
    local display_name="$3"
    local version="${!env_var:-}"

    if [[ -z "${version}" ]]; then
      continue
    fi

    echo "$env_var: $version"
    yq -i ".vars.${yaml_var} = strenv(${env_var})" "$BUILD_TASKFILE"
  done
}

github_labels() {
  # --force also resets the colour and description of an existing label
  gh label create "${DEPENDENCIES_LABEL}" \
    --color "#0366d6" \
    --description "Pull requests that update a dependency file" \
    --force
}

commit_and_push_changes() {
  git add .
  git config user.name "github-actions[bot]"
  git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
  git commit -m "${PR_TITLE}"
  git push origin "${PACKAGE_VERSION_UPDATER_BRANCH}" --force-with-lease
}

create_or_edit_pr() {
  if [[ "$(gh pr list --head "${PACKAGE_VERSION_UPDATER_BRANCH}" --json number --jq length)" != "0" ]]; then
    echo "PR exists already. Updating the 'title' and 'description'..."

    gh pr edit "${PACKAGE_VERSION_UPDATER_BRANCH}" \
      --body "${PR_BODY}" \
      --title "${PR_TITLE}"

    return
  fi

  echo "creating pr..."
  gh pr create \
    --base main \
    --body "${PR_BODY}" \
    --head "${PACKAGE_VERSION_UPDATER_BRANCH}" \
    --label "${DEPENDENCIES_LABEL}" \
    --title "${PR_TITLE}"
}

generate_pr_body_with_updates() {
  local dep
  local pr_body=""

  for dep in "${PACKAGES_TO_BE_UPDATED[@]}"; do
    set -- $dep

    local env_var="$1"
    local yaml_var="$2"
    local display_name="$3"
    local new_version="${!env_var:-}"
    local old_version

    old_version=$(yq -r ".vars.${yaml_var}" "$BUILD_TASKFILE")

    if [[ -z "$new_version" || -z "$old_version" ]]; then
      continue
    fi

    if [[ "$new_version" != "$old_version" ]]; then
      pr_body+="Updated $display_name: $old_version → $new_version"$'\n'
    fi
  done

  export PR_BODY="$pr_body"
  echo "PR_BODY: ${PR_BODY}"
}

update_task_version() {
  local old_version
  local new_version
  local display_name="Task"
  local gitlab_task_version=".spec.inputs.task-version.default"

  old_version=$(yq eval '.inputs.task-version.default' action.yml)
  new_version=$(latest_stable_package_version_on_github go-task/task v3.) || new_version=""

  if [[ -z "${new_version}" || "${new_version}" == "null" ]]; then
    echo "Warning: Failed to fetch latest Task version" >&2
    FAILED_LOOKUPS+=("${display_name}")

    return
  fi

  new_version="${new_version#v}"

  echo "Task version: ${old_version} → ${new_version}"

  if [[ "${old_version}" != "${new_version}" ]]; then
    yq eval ".inputs.task-version.default = \"${new_version}\"" -i action.yml
    yq eval "(select(documentIndex == 0) | ${gitlab_task_version}) = \"${new_version}\"" -i templates/golang.yml

    if [[ $(yq eval '.inputs.task-version.default' action.yml) == "${new_version}" &&
      $(yq eval "select(documentIndex == 0) | ${gitlab_task_version}" templates/golang.yml) == "${new_version}" ]]; then
      PR_BODY+="Updated ${display_name}: ${old_version} → ${new_version}"$'\n'
      echo "Successfully updated Task version"

      return
    fi

    echo "Error: Failed to update Task version in action.yml or templates/golang.yml" >&2

    return 1
  fi

  echo "Task version is already up to date (${old_version})"

  return
}

# GRYPE_CHECKSUMS_SHA256 pins the checksums file of the GRYPE_VERSION release,
# so it has to follow every GRYPE_VERSION bump
update_grype_checksums_sha256() {
  local version
  local checksums_sha256

  version=$(yq -r '.vars.GRYPE_VERSION' "${BUILD_TASKFILE}")
  if ! checksums_sha256=$(curl -fsSL "https://github.com/anchore/grype/releases/download/${version}/grype_${version#v}_checksums.txt" | shasum -a 256 | cut -d ' ' -f 1); then
    echo "Error: failed to download the grype ${version} checksums file" >&2

    return 1
  fi

  yq -i ".vars.GRYPE_CHECKSUMS_SHA256 = \"${checksums_sha256}\"" "${BUILD_TASKFILE}"
}

# pins the golang image of the GitLab CI/CD components to the Go version of
# go.mod and the current digest of that tag, which changes when the image is
# rebuilt, e.g. for OS updates
update_golang_image() {
  local display_name="golang image"
  local image_default=".spec.inputs.image.default"
  local templates=(templates/golang.yml templates/golang-release.yml)
  local go_version
  local token
  local digest
  local old_image
  local new_image
  local template

  go_version=$(awk '$1 == "go" { print $2 }' go.mod)
  token=$(curl -fsS "https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/golang:pull" | jq -r .token) || token=""
  digest=$(curl -fsSI \
    -H "Authorization: Bearer ${token}" \
    -H "Accept: application/vnd.oci.image.index.v1+json" \
    -H "Accept: application/vnd.docker.distribution.manifest.list.v2+json" \
    "https://registry-1.docker.io/v2/library/golang/manifests/${go_version}" | \
      tr -d '\r' | awk 'tolower($1) == "docker-content-digest:" { print $2 }') || digest=""

  if [[ -z "${go_version}" || "${digest}" != sha256:* ]]; then
    echo "Warning: failed to determine the golang image digest" >&2
    FAILED_LOOKUPS+=("${display_name}")

    return
  fi

  old_image=$(yq -r "select(documentIndex == 0) | ${image_default}" "${templates[0]}")
  new_image="golang:${go_version}@${digest}"

  if [[ "${old_image}" == "${new_image}" ]]; then
    echo "golang image is already up to date (${old_image})"

    return
  fi

  for template in "${templates[@]}"; do
    yq -i "(select(documentIndex == 0) | ${image_default}) = \"${new_image}\"" "${template}"
  done

  PR_BODY+="Updated ${display_name}: ${old_image} → ${new_image}"$'\n'
}

main() {
  set -xeuo pipefail

  latest_stable_package_versions
  checkout_branch_required_to_apply_package_version_updates
  generate_pr_body_with_updates
  update_task_version
  update_golang_image
  replace_versions_with_latest_stable_package_versions
  update_grype_checksums_sha256

  if [[ -n "$(git status --porcelain)" ]]; then
    github_labels
    commit_and_push_changes
    create_or_edit_pr
  else
    echo "No changes to commit."
  fi

  if [[ ${#FAILED_LOOKUPS[@]} -gt 0 ]]; then
    echo "Error: failed to determine a version for: ${FAILED_LOOKUPS[*]}" >&2

    return 1
  fi
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  main "$@"
fi
