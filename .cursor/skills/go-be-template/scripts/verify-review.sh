#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "$0")" && pwd)/verify-lib.sh"
DIR="$(cd "$(dirname "$0")" && pwd)"

"$DIR/verify-structure.sh"
"$DIR/verify-openapi.sh"
"$DIR/verify-data.sh"

ok "go-be-template-review"
