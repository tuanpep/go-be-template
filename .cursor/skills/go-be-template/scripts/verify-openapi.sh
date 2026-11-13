#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

[[ -f api/openapi.yaml ]] || fail "missing api/openapi.yaml"
grep -q 'operationId:' api/openapi.yaml || fail "openapi.yaml must define operationId on operations"

if go tool oapi-codegen -version >/dev/null 2>&1; then
  go tool oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml >/dev/null
elif command -v oapi-codegen >/dev/null; then
  oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml >/dev/null
else
  fail "oapi-codegen missing — use the go.mod tool (go tool oapi-codegen)"
fi

if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  git diff --exit-code internal/transport/openapi || fail "api.gen.go drift — run make generate"
fi

test -f internal/transport/api.go || fail "missing transport/api.go"
grep -q 'ServerInterface' internal/transport/api.go || fail "transport.API must implement ServerInterface"

ok "go-be-template-openapi"
