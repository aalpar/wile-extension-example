#!/usr/bin/env bash
# Smoke test: verify built example binaries exist and are executable.
#
# Usage:
#   ./tools/sh/smoke-test.sh <dist-dir>

set -euo pipefail

DIST_DIR="${1:?Usage: smoke-test.sh <dist-dir>}"

if [ ! -d "$DIST_DIR" ]; then
    echo "Error: dist directory not found: $DIST_DIR"
    exit 1
fi

passed=0
failed=0
total=0

for bin in "$DIST_DIR"/*; do
    [ -f "$bin" ] || continue
    name=$(basename "$bin")
    total=$((total + 1))

    echo -n "  $name executable... "
    if [ -x "$bin" ]; then
        echo "ok"
        passed=$((passed + 1))
    else
        echo "FAIL (not executable)"
        failed=$((failed + 1))
    fi
done

if [ "$total" -eq 0 ]; then
    echo "FAIL: no binaries found in $DIST_DIR"
    exit 1
fi

echo ""
echo "Smoke tests: $passed/$total passed"
if [ "$failed" -gt 0 ]; then
    exit 1
fi
