#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

[[ -f internal/platform/log/log.go ]] || fail "missing platform/log"

for h in internal/*/handler.go; do
  [[ -f "$h" ]] || continue
  grep -q 'internal/platform/httpx' "$h" || fail "$h must import platform/httpx for Problem responses"
done

ok "go-be-template-errors"
