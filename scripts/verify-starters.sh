#!/usr/bin/env bash
# Checks that every starter's tests pass once the reference solution replaces the stubs.
# Usage: scripts/verify-starters.sh [java|go]...   (default: both)
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

verify_java() {
  cp -r "$root/java" "$work/java"
  rm -rf "$work/java/build" "$work/java/.gradle"
  (cd "$root/java" && find src/main/java -name '*.java') | while read -r file; do
    cp "$root/solutions/java/$file" "$work/java/$file"
  done
  (cd "$work/java" && ./gradlew test --console=plain -q)
  echo "java: starter tests pass with the reference solution"
}

verify_go() {
  cp -r "$root/go" "$work/go"
  (cd "$root/go" && find . -name '*.go' ! -name '*_test.go' ! -path './algorithms/*') | while read -r file; do
    sed 's#labs/solutions/go#labs/go#' "$root/solutions/go/$file" > "$work/go/$file"
  done
  (cd "$work/go" && go vet ./... && go test -race -count=1 ./...)
  echo "go: starter tests pass with the reference solution"
}

for language in "${@:-java go}"; do
  for l in $language; do "verify_$l"; done
done
