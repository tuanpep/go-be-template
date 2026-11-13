#!/usr/bin/env bash
set -euo pipefail
# shellcheck source=verify-lib.sh
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
verify_repo_root

assert_services_no_store_db
assert_handlers_no_store_db
assert_services_no_openapi
assert_platform_no_features

go build -o /dev/null ./...

ok "go-be-template-structure"
