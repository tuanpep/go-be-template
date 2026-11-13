#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

[[ -f internal/app/app.go ]] || fail "missing internal/app/app.go"
grep -q 'transport.NewAPI' internal/app/app.go || fail "app must construct transport.API"
grep -q 'httpserver.NewRouter' internal/app/app.go || fail "app must construct router"

lines=$(wc -l < cmd/api/main.go)
if [[ "$lines" -gt 80 ]]; then
  fail "cmd/api/main.go too large ($lines lines) — move wiring to internal/app"
fi

if grep -q 'chi\.NewRouter\|HandlerFromMux\|HandlerWithOptions' cmd/api/main.go 2>/dev/null; then
  fail "route registration belongs in internal/http and internal/app, not cmd/api"
fi

go build -o /dev/null ./cmd/api ./cmd/migrate

ok "go-be-template-wiring"
