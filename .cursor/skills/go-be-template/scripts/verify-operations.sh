#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

[[ -f .env.example ]] || fail "missing .env.example"
[[ -f docker-compose.yml ]] || fail "missing docker-compose.yml"
[[ -f Dockerfile ]] || fail "missing Dockerfile"
[[ -f cmd/migrate/main.go ]] || fail "missing cmd/migrate"

grep -q 'HTTP_ADDR' .env.example || fail ".env.example should document HTTP_ADDR"
grep -q 'DATABASE_URL' .env.example || fail ".env.example should document DATABASE_URL"

go build -o /dev/null ./cmd/migrate

ok "go-be-template-operations"
