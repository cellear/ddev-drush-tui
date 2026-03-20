# Sprint 3: Command Execution

**Sprint Goal:** Pressing Run in the parameter form executes the Drush command and displays real output in the bottom pane. The TUI doesn't freeze during execution.

**Test site:** `drupal-cms/` in the repo root (gitignored; must be running via `ddev start`)

---

## Stories

### S3-1 · Command executor · [x]

- `internal/drush/executor.go`: build command from form values, run via `exec.Command`, capture combined stdout/stderr, 60s timeout, proper exit code handling

### S3-2 · Output pane · [x]

- `internal/tui/output.go`: scrollable `tview.TextView` with "Running..." state, result display, red-highlighted errors for non-zero exit codes

### S3-3 · Wire Run button · [x]

- `internal/tui/params.go`: Run button collects positional args and `--option=value` flags from form fields
- `internal/tui/app.go`: wires Run callback to executor in a goroutine with `QueueUpdateDraw`

### S3-4 · Form navigation fix · [x]

- Arrow down/up now moves between form fields (translated to Tab/Backtab)
- Tab from command list enters the form; Tab within the form navigates fields
- Fixed deadlock caused by calling `app.Draw()` inside an event handler

### S3-5 · Demo scripts · [x]

- `scripts/demo-s3.sh` with Sprint 3 checklist
- All demo scripts now pause for Enter before launching the TUI

---

## Acceptance

**Status:** Accepted
**Date:** 2026-03-20
**Reviewed by:** Luke McCormick

> THAT WORKED! That was actually my first actual output from our program. It's a real thing :) :) :)
> Pretty good. Still things to work on, but it's already good enough to demo.

---

## Decisions Made This Sprint

- **Demo scripts pause for Enter:** All demo scripts now wait for the user to press Enter before launching the TUI, so you can read the checklist first.
- **Arrow keys navigate forms:** Down/Up arrow keys move between form fields, matching the command list navigation pattern.

---

## Deferred to Future Sprint

- Output pane scrolling (output is scrollable in code but needs a way for the user to focus/scroll it)
- Command list search/filter (`/` key)
- Required argument validation before Run

---

## References

- `AGENT.md` — project constraints, UI layout, development process
- `DOC/implementation-plan.md` — Phase 8 and 9 details
- `SPRINTS/sprint-2.md` — previous sprint
