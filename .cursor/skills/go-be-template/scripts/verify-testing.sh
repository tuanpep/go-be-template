#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

[[ -f .github/workflows/ci.yml ]] || fail "missing CI workflow"
grep -q 'git diff' .github/workflows/ci.yml || fail "CI should check generate drift"
grep -q 'go test -race' .github/workflows/ci.yml || fail "CI should run go test -race"
grep -q 'govulncheck' .github/workflows/ci.yml || fail "CI should run govulncheck"

go test ./...

ok "go-be-template-testing"
