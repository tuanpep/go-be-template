#!/usr/bin/env bash
# Build, test, generate drift, and section checks. Run from anywhere.
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$DIR/../../../.." && pwd)"
cd "$ROOT"

for script in verify-structure verify-openapi verify-data verify-http verify-wiring verify-errors verify-testing verify-operations verify-security; do
  bash "$DIR/$script.sh"
done

generated=0
if go tool sqlc version >/dev/null 2>&1 && go tool oapi-codegen -version >/dev/null 2>&1; then
  go tool sqlc generate
  go tool oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
  generated=1
elif command -v sqlc >/dev/null && command -v oapi-codegen >/dev/null; then
  sqlc generate
  oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
  generated=1
else
  echo "warn: sqlc/oapi-codegen not available; skipping generate drift check"
fi
if [[ "$generated" -eq 1 ]] && git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  git diff --exit-code internal/store/db internal/transport/openapi || {
    echo "generated code out of date — run make generate"
    exit 1
  }
fi

go build -o /dev/null ./cmd/api ./cmd/migrate
go test ./...

if command -v golangci-lint >/dev/null; then
  golangci-lint run ./...
else
  echo "warn: golangci-lint not installed; skipping"
fi

echo "OK: go-be-template validation passed"
