#!/bin/bash
set -e

# Sprint 1 Demo — run this from the drupal-cms/ directory.
#
# PURPOSE: This is a DEMO script, not an automated test. It runs the
# program and shows the output so a human can visually verify that
# the sprint's work is correct. The human decides pass/fail.
#
# Sprint: 1
# Stories completed: S1-1 (scaffold), S1-2 (DDEV detection), S1-3 (discovery)
# Changes:
#   - internal/ddev/context.go: Detects DDEV project root and project name.
#   - internal/drush/discovery.go: Parses 'ddev drush list --format=json' into Go structs.
#   - cmd/ddev-drush-tui/main.go: Entry point wiring detection + discovery.

echo "===================================================="
echo "   Sprint 1 Demo: Command Discovery                 "
echo "===================================================="
echo ""
echo "WHAT TO LOOK FOR:"
echo "1. Project detected (e.g. 'drupal-cms-drush-tui')"
echo "2. Drush commands listed, grouped by namespace"
echo "3. No errors or panics"
echo "----------------------------------------------------"

# Run the program — human reviews the output below.
go run ../cmd/ddev-drush-tui/main.go

echo "----------------------------------------------------"
echo "Demo complete. Review the output above."
echo "===================================================="
