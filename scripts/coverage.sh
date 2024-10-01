#!/usr/bin/env bash

shopt -s extglob
set -e

rnd(){
    echo $(cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f0-9' | head -c 32)
}
fcomp() {
    awk "BEGIN {exit !($1+0<$2+0)}"
}
COVERAGE_DIR="./out/var/$(rnd)"

mkdir -p "${COVERAGE_DIR}"
if [[ $? -ne 0 ]] ; then
    echo "failed"
    exit 1
fi
COVERAGE_DIR="$(realpath ${COVERAGE_DIR})"
COVERAGE_FILE="${COVERAGE_DIR}"/coverage.cov
# count, atomic
COVER_MODE="set"
PKG_LIST=$(go list ./...)

for package in ${PKG_LIST}; do
    go test -covermode="$COVER_MODE" -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package" >/dev/null
done

# Merge the coverage profile files
FILE_LIST=$(ls "$COVERAGE_DIR"/*.cov)
echo "mode: ${COVER_MODE}" >"${COVERAGE_FILE}"
for file in ${FILE_LIST}; do
  echo "copying contents of ${file} into coverage.cov"
  tail -q -n +2 "${file}" >>"${COVERAGE_FILE}"
done

# Display the global code coverage
coveragePercentage=$(go tool cover -func="${COVERAGE_FILE}" | grep -E '^total:' | awk '{print $3}')
coveragePercentage=${coveragePercentage%\%}
echo "${coveragePercentage}" > ./out/var/coverage
go tool cover -html "${COVERAGE_FILE}" -o ./out/var/coverage.html

echo "${coveragePercentage}"