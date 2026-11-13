#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

if grep -rn 'fmt\.Sprintf.*SELECT\|fmt\.Sprintf.*INSERT' internal/ --include='*.go' 2>/dev/null; then
  fail "use sqlc parameterized queries, not fmt.Sprintf SQL"
fi

if grep -qE 'password|secret' api/openapi.yaml 2>/dev/null; then
  echo "warn: check openapi.yaml for accidental secret examples"
fi

[[ ! -f .env ]] || fail ".env must not be committed (use .env.example)"

ok "go-be-template-security"
