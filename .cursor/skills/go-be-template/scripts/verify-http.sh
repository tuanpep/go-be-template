#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

[[ -f internal/http/router.go ]] || fail "missing internal/http/router.go"
grep -q 'HandlerWithOptions' internal/http/router.go || fail "router must register OpenAPI via HandlerWithOptions"
grep -q 'ErrorHandlerFunc' internal/http/router.go || fail "parameter errors must use the OpenAPI ErrorHandlerFunc"
grep -q 'RequestID' internal/http/router.go || fail "router must use RequestID middleware"

[[ -f internal/platform/httpx/json.go ]] || fail "missing platform/httpx/json.go"
[[ -f internal/platform/httpx/problem.go ]] || fail "missing platform/httpx/problem.go"

ok "go-be-template-http"
