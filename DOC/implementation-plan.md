# Implementation Plan

**Project:** ddev-drush-tui
**Stack:** Go + tview
**Module:** github.com/cellear/ddev-drush-tui
**DDEV command:** `ddev drush-tui`

---

## Architecture

```
cmd/ddev-drush-tui/
  main.go                 # Entry point: detect DDEV project, init TUI

internal/
  ddev/
    context.go            # Detect DDEV project root, project name, verify drush

  drush/
    discovery.go          # ddev drush list --format=json → []Command structs
    help.go               # ddev drush help <cmd> --format=json → CommandHelp struct
    executor.go           # Build + run ddev drush command, capture output

  tui/
    app.go                # tview Application, layout, focus management, keybindings
    commands.go           # Left pane: scrollable list, grouped by namespace
    params.go             # Right pane: dynamic form for selected command's args/options
    output.go             # Bottom pane: output from last command run

commands/host/drush-tui   # DDEV add-on host command script
install.yaml              # DDEV add-on manifest
Makefile                  # build, install, dist targets
```

---

## Phase Breakdown

### Phase 1: Scaffold (Sprint 1, S1-1)

**Owner:** Codex (alt: Haiku)
**Goal:** `go build ./...` succeeds with zero errors.

- `go mod init github.com/cellear/ddev-drush-tui`
- Add tview + tcell dependencies
- Create all package directories with stub files
- Stub `main.go` that prints "ddev-drush-tui" and exits

Nothing functional yet. Just a compiling Go module.

---

### Phase 2: DDEV Context Detection (Sprint 1, S1-2)

**Owner:** Codex (alt: Gemini CLI)
**Goal:** Know which DDEV project we're in, and whether drush is available.

`internal/ddev/context.go`:
- Run `ddev describe --json` to get project name
- Verify `ddev drush version` exits 0
- Return a `Context` struct: `{ProjectName string, AppRoot string}`
- Exit with clear error if not in a DDEV project or drush unavailable

No TUI yet. Wire to `main.go` to print project name.

**Reference:** See `DOC/drush-discovery.md` for how to invoke ddev commands.

---

### Phase 3: Command Discovery (Sprint 1, S1-3)

**Owner:** Gemini CLI (alt: Sonnet)
**Goal:** Parse `ddev drush list --format=json` into Go structs.

`internal/drush/discovery.go`:
- Run `ddev drush list --format=json`, capture stdout
- Unmarshal JSON into `CommandList` struct (see `DOC/drush-discovery.md` for struct design)
- **Handle the inconsistent `arguments` field** (empty array vs object map) — critical parsing issue
- Extract namespace from command name (everything before `:`)
- Return `[]Command` grouped by namespace
- Cache result in memory for the session

**IMPORTANT:** The `definition.arguments` field in list output uses `[]` for no args but `{}` for args. Use `json.RawMessage` to handle both. See `DOC/drush-discovery.md` for the pattern.

---

### Phase 4: Minimal TUI with Command List (Sprint 1, S1-4)

**Owner:** Cursor Composer (alt: Haiku)
**Goal:** TUI launches showing grouped list of real Drush commands.

`internal/tui/app.go`:
- tview `Application` with three-pane layout
- Header bar showing project name
- Left pane (60%): command list
- Right pane (40%): placeholder "Select a command"
- Bottom pane (25%): placeholder "Output will appear here"
- `q` quits

`internal/tui/commands.go`:
- Populate `tview.List` with commands grouped by namespace
- Namespace headers shown as non-selectable items (grayed)
- Arrow keys navigate, Enter selects (wires to params pane in Sprint 2)

**Sprint 1 Demo checkpoint:**
`ddev drush-tui` (or `go run ./cmd/ddev-drush-tui`) shows TUI with real Drush commands, grouped by namespace, navigable with arrow keys, `q` exits.

---

### Phase 5: DDEV Add-on Stub (Sprint 1, S1-5)

**Owner:** Codex (alt: Haiku)
**Goal:** `ddev drush-tui` launches the TUI.

- `install.yaml` listing the host command
- `commands/host/drush-tui` shell script (requires `#ddev-generated` as line 2)
- `Makefile` with `build`, `install`, `dist` targets
- Binary name: `ddev-drush-tui`

**CRITICAL DDEV add-on convention:** Source files are stored WITHOUT a `.ddev/` prefix in the repo. DDEV prepends it on install. See `XDEBUG-SAVED/DOC/ddev-addon-conventions.md` for the full rules from the xdebug project.

---

### Phase 6: Parameter Discovery (Sprint 2, S2-1)

**Owner:** Codex (alt: Sonnet)
**Goal:** Selecting a command fetches and caches its parameters.

`internal/drush/help.go`:
- Run `ddev drush help <command> --format=json`, capture stdout
- Unmarshal into `CommandHelp` struct (see `DOC/drush-discovery.md`)
- Filter out global Drush options (druplicon, notify, xh-link, etc.) — see filter list in discovery doc
- Cache per command name for the session

---

### Phase 7: Parameter Form (Sprint 2, S2-2 + S2-3)

**Owner:** Cursor Composer (alt: Haiku)
**Goal:** Right pane shows a tview Form for the selected command.

`internal/tui/params.go`:
- When command selected: call `drush.Help()`, render form
- Required arguments: `tview.InputField` labeled with `*` prefix
- Optional string options: `tview.InputField`
- Boolean flags (accept_value="0"): `tview.Checkbox`
- Options with defaults: pre-filled
- "Run" button at bottom (disabled until required args filled)
- "Cancel" button returns focus to command list

Form validation: required arguments must not be empty. Show error in status bar.

---

### Phase 8: Command Execution (Sprint 3, S3-1 + S3-2)

**Owner:** Codex (alt: Haiku)
**Goal:** Run command, show output in bottom pane.

`internal/drush/executor.go`:
- Build command: `ddev drush <name> [arg1] [arg2] [--opt1=val] [--flag]`
- Run via `exec.Command`, capture stdout+stderr
- Return output string and exit code

`internal/tui/output.go`:
- Show command header: `$ ddev drush cache:rebuild`
- Stream or display output in `tview.TextView`
- Scrollable
- Error exits highlighted (non-zero exit code)

Run command in goroutine to avoid freezing TUI. Use `app.QueueUpdateDraw()` for thread-safe UI updates.

**Sprint 3 Demo checkpoint:**
Select `cache:rebuild`, hit Run → output pane shows real Drush output.
Select `core:status` → shows site status table.
Select a failing command → shows error output clearly.

---

### Phase 9: Error Handling (Sprint 3, S3-3)

**Owner:** Sonnet
**Goal:** Graceful handling of all failure modes.

- DDEV not running: detect before launch, print clear message
- Drush not available: detect before launch
- Command timeout: configurable, default 60s, show timeout message
- DDEV stops mid-session: handle subprocess error on execute
- Unknown command: shouldn't happen (we built the list from drush), but guard anyway

---

### Phase 10: Polish + Release (Sprint 4)

**Owner:** Sonnet (architecture), Haiku/Codex (implementation)

- `/` to filter command list
- Command description shown in right pane before form is loaded
- Command aliases shown
- Cross-platform build: darwin/arm64, darwin/amd64, linux/amd64
- GitHub release workflow
- README.md with usage, installation, screenshots

---

## Key Constraints (from AGENT.md)

- No persistent config (PoC scope)
- No multi-site handling
- No auto-execution without user confirmation
- Single binary, cross-compiled
- tview for UI (no web UI, no Bubble Tea, etc.)
- All Drush calls via `ddev drush` — never call drush directly

---

## Known Implementation Challenges

1. **Inconsistent arguments JSON** — The `list` output uses `[]` vs `{}` for arguments. Must handle with `json.RawMessage`. Documented in `DOC/drush-discovery.md`.

2. **Global option filtering** — Every command inherits ~20 global Drush options. Filter list in `DOC/drush-discovery.md`.

3. **Long-running commands** — `drush deploy`, `drush cr` on large sites can take 30+ seconds. Run in goroutine with spinner in output pane.

4. **Commands requiring interactive input** — Some commands (`drush pm:install` with prompts) are out of PoC scope. Detect and warn.

5. **tview focus management** — Three panes with Tab cycling. Make sure Run button in form can be triggered with Enter, not just mouse.

---

*Last updated: 2026-03-19 by claude-sonnet*
