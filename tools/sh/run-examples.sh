#!/usr/bin/env bash
# Run all example binaries from a dist directory and verify they exit 0.
# Each example gets a 30-second timeout to catch hangs.
#
# Usage:
#   ./tools/sh/run-examples.sh <dist-dir>

set -euo pipefail

DIST_DIR="${1:?Usage: run-examples.sh <dist-dir>}"

if [ ! -d "$DIST_DIR" ]; then
    echo "Error: dist directory not found: $DIST_DIR"
    exit 1
fi

TIMEOUT="${EXAMPLE_TIMEOUT:-30}"

# Detect timeout command (GNU timeout or macOS gtimeout via coreutils)
if command -v timeout >/dev/null 2>&1; then
    TIMEOUT_CMD=(timeout --kill-after=5)
elif command -v gtimeout >/dev/null 2>&1; then
    TIMEOUT_CMD=(gtimeout --kill-after=5)
else
    TIMEOUT_CMD=()
fi

passed=0
failed=0
failures=()

for bin in "$DIST_DIR"/*; do
    [ -x "$bin" ] || continue
    [ -f "$bin" ] || continue

    name=$(basename "$bin")
    echo -n "  $name... "

    if "${TIMEOUT_CMD[@]}" "$TIMEOUT" "$bin" >/dev/null 2>&1; then
        echo "ok"
        passed=$((passed + 1))
    else
        echo "FAIL"
        failed=$((failed + 1))
        failures+=("$name")
    fi
done

total=$((passed + failed))
echo ""
echo "Examples: $passed passed, $failed failed ($total total)"

if [ "$failed" -gt 0 ]; then
    echo ""
    echo "Failures:"
    for f in "${failures[@]}"; do
        echo "  FAIL  $f"
    done
    exit 1
fi
