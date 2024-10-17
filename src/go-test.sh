#!/bin/bash

COVERPKG=$1
COVERPROFILE=$2
GOLANG_PARALLEL_TESTS=$3
GOLANG_TEST_EXCLUSIONS=$4
TAGS=$5

variables() {
  if [ -z "${GOLANG_PARALLEL_TESTS}" ]; then
    if [ "$(uname -s)" = "Darwin" ]; then
      GOLANG_PARALLEL_TESTS=$(sysctl -n hw.ncpu)
    else
      GOLANG_PARALLEL_TESTS=$(nproc)
    fi
  fi

  if [ -n "${COVERPROFILE}" ]; then
    COVERPROFILE="-coverprofile=${COVERPROFILE}"
  fi

  if [ -n "${TAGS}" ]; then
    TAGS="--tags=${TAGS}"
  fi

  GO_LIST_TEST_EXCLUSIONS=$(go list $TAGS ./... | grep -v ${GOLANG_TEST_EXCLUSIONS})

  if [ -n "${COVERPKG}" ]; then
    COVERPKG="-coverpkg=$(echo "${GO_LIST_TEST_EXCLUSIONS}" | tr '\n' ',')"
  fi

  echo "COVERPKG: ${COVERPKG}"
  echo "COVERPROFILE: ${COVERPROFILE}"
  echo "GO_LIST_TEST_EXCLUSIONS: ${GO_LIST_TEST_EXCLUSIONS}"
  echo "GOLANG_PARALLEL_TESTS: ${GOLANG_PARALLEL_TESTS}"
  echo "GOLANG_TEST_EXCLUSIONS: ${GOLANG_TEST_EXCLUSIONS}"
  echo "TAGS: ${TAGS}"
}

run_go_tests() {
  go test \
    -p "${GOLANG_PARALLEL_TESTS}" \
    -race \
    -short \
    -v \
    $COVERPKG \
    $COVERPROFILE \
    $TAGS \
    $(echo "${GO_LIST_TEST_EXCLUSIONS}")
}

main() {
  variables
  run_go_tests
}

main
