#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

[[ -d migrations ]] || fail "missing migrations/"
[[ -f sql/schema/schema.sql ]] || fail "missing sql/schema/schema.sql"
[[ -d sql/queries ]] || fail "missing sql/queries/"

assert_services_no_store_db

generated=0
if go tool sqlc version >/dev/null 2>&1; then
  go tool sqlc generate
  generated=1
elif command -v sqlc >/dev/null; then
  sqlc generate
  generated=1
else
  echo "warn: sqlc not available; skipping generate drift"
fi
if [[ "$generated" -eq 1 ]] && git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  git diff --exit-code internal/store/db || fail "sqlc drift — run make generate"
fi

if grep -q 'goose\.\|migrate\.Up' cmd/api/main.go 2>/dev/null; then
  fail "cmd/api must not run migrations"
fi

ok "go-be-template-data"
