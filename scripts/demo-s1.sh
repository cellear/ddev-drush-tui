#!/bin/bash
set -e

# Sprint 1 Demo — run from the repo root: ./scripts/demo-s1.sh
#
# PURPOSE: This is a DEMO script, not an automated test. It runs the
# program and shows the output so a human can visually verify that
# the sprint's work is correct. The human decides pass/fail.
#
# REQUIRES: The drupal-cms/ DDEV test site must be running (ddev start).
#
# Sprint: 1
# Stories completed: S1-1 (scaffold), S1-2 (DDEV detection), S1-3 (discovery)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_SITE="$REPO_ROOT/drupal-cms"

if [ ! -d "$TEST_SITE" ]; then
  echo "Error: drupal-cms/ test site not found at $TEST_SITE"
  exit 1
fi

echo "===================================================="
echo "   Sprint 1 Demo: Command Discovery                 "
echo "===================================================="
echo ""
echo "WHAT TO LOOK FOR:"
echo "1. Project detected (e.g. 'drupal-cms-drush-tui')"
echo "2. Drush commands listed, grouped by namespace"
echo "3. No errors or panics"
echo "----------------------------------------------------"

# Run the program from the test site directory so DDEV context is detected.
cd "$TEST_SITE"
go run "$REPO_ROOT/cmd/ddev-drush-tui/main.go"

echo "----------------------------------------------------"
echo "Demo complete. Review the output above."
echo "===================================================="
