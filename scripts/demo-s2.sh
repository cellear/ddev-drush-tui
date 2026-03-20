#!/bin/bash
set -e

# Sprint 2 Demo — run from the repo root: ./scripts/demo-s2.sh
#
# PURPOSE: This is a DEMO script, not an automated test. It runs the
# program and shows the output so a human can visually verify that
# the sprint's work is correct. The human decides pass/fail.
#
# REQUIRES: The drupal-cms/ DDEV test site must be running (ddev start).
#
# Sprint: 2
# Stories: S2-1 (help parsing), S2-2 (wiring), S2-3 (param form)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_SITE="$REPO_ROOT/drupal-cms"

if [ ! -d "$TEST_SITE" ]; then
  echo "Error: drupal-cms/ test site not found at $TEST_SITE"
  exit 1
fi

echo "===================================================="
echo "   Sprint 2 Demo: Parameter Discovery + Form         "
echo "===================================================="
echo ""
echo "WHAT TO LOOK FOR:"
echo "1. TUI launches with grouped command list (Sprint 1)"
echo "2. Select user:create → form shows '* name' (required)"
echo "3. Select cache:rebuild → form shows options, no globals"
echo "4. Tab cycles focus between command list and form"
echo "5. Esc returns from form to command list"
echo "6. Cancel button returns to command list"
echo "7. q exits from command list"
echo "----------------------------------------------------"
echo ""
read -rp "Press Enter to launch the TUI..."

# Run the program from the test site directory.
cd "$TEST_SITE"
go run "$REPO_ROOT/cmd/ddev-drush-tui/main.go"

echo "----------------------------------------------------"
echo "Demo complete. Review the output above."
echo "===================================================="
