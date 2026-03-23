#!/bin/bash
set -e

# Sprint 4 Demo — run from the repo root: ./scripts/demo-s4.sh
#
# PURPOSE: This is a DEMO script, not an automated test. It runs the
# program and shows the output so a human can visually verify that
# the sprint's work is correct. The human decides pass/fail.
#
# REQUIRES: The drupal-cms/ DDEV test site must be running (ddev start).
#
# Sprint: 4
# Stories: S4-1 (hierarchical commands), S4-2 (output focus/scroll),
#          S4-3 (filter), S4-4 (required-arg validation), S4-5 (demo)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TEST_SITE="$REPO_ROOT/drupal-cms"

if [ ! -d "$TEST_SITE" ]; then
  echo "Error: drupal-cms/ test site not found at $TEST_SITE"
  exit 1
fi

echo "===================================================="
echo "   Sprint 4 Demo: Navigation + UX Polish            "
echo "===================================================="
echo ""
echo "WHAT TO LOOK FOR:"
echo "1. TUI launches showing namespace list, not a flat command list"
echo "2. Enter on a namespace drills into its commands"
echo "3. ← Back and Esc return to the namespace list"
echo "4. Select a command, fill params, Run -> output appears"
echo "5. Tab cycles commands -> params -> output -> commands"
echo "6. Output border changes color when the output pane is focused"
echo "7. Output scroll keys work: arrows, Space, b, Home/g, End/G"
echo "8. / opens the inline filter; typing narrows the list; Esc clears it"
echo "9. Run user:create with empty name -> red validation error appears"
echo "10. Typing in the required field clears the validation error"
echo "11. q exits from the command list"
echo "----------------------------------------------------"
echo ""
echo "Suggested manual path:"
echo "- Drill into cache, then use ← Back"
echo "- Open cache:rebuild and run it"
echo "- Tab to output and test scrolling"
echo "- Use / to filter namespaces or commands"
echo "- Open user:create and press Run before filling name"
echo "----------------------------------------------------"
echo ""
read -rp "Press Enter to launch the TUI..."

# Run the program from the test site directory so DDEV context is detected.
cd "$TEST_SITE"
go run "$REPO_ROOT/cmd/ddev-drush-tui/main.go"

echo "----------------------------------------------------"
echo "Demo complete. Review the behavior above."
echo "===================================================="
