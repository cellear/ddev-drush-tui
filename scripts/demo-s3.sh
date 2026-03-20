#!/bin/bash
set -e

# Sprint 3 Demo — run from the repo root: ./scripts/demo-s3.sh
#
# PURPOSE: This is a DEMO script, not an automated test. It runs the
# program and shows the output so a human can visually verify that
# the sprint's work is correct. The human decides pass/fail.
#
# REQUIRES: The drupal-cms/ DDEV test site must be running (ddev start).
#
# Sprint: 3
# Stories: S3-1 (executor), S3-2 (output pane), S3-3 (wiring), S3-4 (goroutine)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_SITE="$REPO_ROOT/drupal-cms"

if [ ! -d "$TEST_SITE" ]; then
  echo "Error: drupal-cms/ test site not found at $TEST_SITE"
  exit 1
fi

echo "===================================================="
echo "   Sprint 3 Demo: Command Execution                  "
echo "===================================================="
echo ""
echo "WHAT TO LOOK FOR:"
echo "1. Select cache:rebuild, press Run → output pane shows real Drush output"
echo "2. Select core:status, press Run → shows site status"
echo "3. Output pane shows '$ ddev drush <command>' header"
echo "4. While command runs, output shows 'Running...'"
echo "5. Failed commands show exit code in red"
echo "6. Output is scrollable"
echo "7. TUI doesn't freeze during command execution"
echo "----------------------------------------------------"
echo ""
read -rp "Press Enter to launch the TUI..."

# Run the program from the test site directory.
cd "$TEST_SITE"
go run "$REPO_ROOT/cmd/ddev-drush-tui/main.go"

echo "----------------------------------------------------"
echo "Demo complete. Review the output above."
echo "===================================================="
