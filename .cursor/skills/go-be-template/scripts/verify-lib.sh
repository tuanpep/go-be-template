# shellcheck shell=bash
# Shared helpers. Source from verify-*.sh. Do not execute directly.

verify_repo_root() {
  local root
  root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)"
  cd "$root"
}

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

ok() {
  echo "OK: $1"
}

assert_services_no_store_db() {
  local bad
  bad="$(find internal -name service.go -print0 | xargs -0 grep -l 'internal/store/db' 2>/dev/null || true)"
  if [[ -n "$bad" ]]; then
    fail "service must not import internal/store/db — use Repository: $bad"
  fi
}

assert_handlers_no_store_db() {
  local bad
  bad="$(find internal -name handler.go -print0 | xargs -0 grep -l 'internal/store/db' 2>/dev/null || true)"
  if [[ -n "$bad" ]]; then
    fail "handler must not import internal/store/db: $bad"
  fi
}

assert_services_no_openapi() {
  local bad
  bad="$(find internal -name service.go -print0 | xargs -0 grep -l 'internal/transport/openapi' 2>/dev/null || true)"
  if [[ -n "$bad" ]]; then
    fail "service must not import OpenAPI types: $bad"
  fi
}

assert_platform_no_features() {
  local bad=""
  for pkg in internal/platform/*/; do
    if grep -rq 'github.com/tuanbt/go-be-template/internal/user' "$pkg" 2>/dev/null; then
      bad="$pkg"
    fi
  done
  if [[ -n "$bad" ]]; then
    fail "platform must not import feature packages: $bad"
  fi
}
