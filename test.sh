#!/usr/bin/env bash
# Run all tests.
#
# Usage:
#   ./test.sh           — Go unit + WS integration + JS unit
#   ./test.sh --e2e     — adds Playwright E2E (requires app running at localhost:52674)
#   APP_URL=https://ohmycode.work ./test.sh --e2e

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")" && pwd)"
RUN_E2E=false
FAIL=0

for arg in "$@"; do
  case "$arg" in
    --e2e) RUN_E2E=true ;;
    *) echo "Unknown argument: $arg"; exit 1 ;;
  esac
done

pass() { echo "  ok  $1"; }
fail() { echo "  FAIL  $1"; FAIL=1; }

# ── Go ──────────────────────────────────────────────────────────────────────
echo "=== Go tests (api) ==="
if (cd "$REPO_ROOT/api" && go test -race ./...); then
  pass "api"
else
  fail "api"
fi

# ── JS unit ─────────────────────────────────────────────────────────────────
echo ""
echo "=== JS unit tests ==="
if (cd "$REPO_ROOT/api/internal/api/client" && npm test --silent); then
  pass "js-unit"
else
  fail "js-unit"
fi

# ── E2E ─────────────────────────────────────────────────────────────────────
if $RUN_E2E; then
  echo ""
  echo "=== E2E tests (Playwright, ${APP_URL:-http://localhost:52674}) ==="
  if (cd "$REPO_ROOT/e2e" && npm install --silent && npx playwright install chromium --quiet && npx playwright test); then
    pass "e2e"
  else
    fail "e2e"
  fi
fi

echo ""
if [ $FAIL -ne 0 ]; then
  echo "FAIL"
  exit 1
fi
echo "PASS"
