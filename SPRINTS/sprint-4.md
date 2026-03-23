# Sprint 4: Navigation + UX Polish

**Sprint Goal:** The command list uses drill-down navigation (namespace → commands) instead of a flat scrollable list. The output pane is focusable and scrollable. `/` filters the list. Required arguments are validated before Run.

**Test site:** `drupal-cms/` in the repo root (gitignored; must be running via `ddev start`)

---

## Stories

### S4-1 · Hierarchical command browsing · [ ]

**Owner:** Cursor Composer (alt: Sonnet)

**Scope:** `internal/tui/commands.go`, `internal/tui/app.go`

Replace the flat grouped list with a two-level drill-down navigation.

**Level 1 — Namespace list:**

```
Commands
─────────────────
  cache          (6)
  checklistapi   (2)
  config         (7)
  core           (5)
  deploy         (4)
  ...
  other          (5)
```

- Each item shows the namespace name and command count
- Enter drills into the selected namespace
- `q` quits

**Level 2 — Commands within namespace:**

```
Commands > cache
─────────────────
  ← Back
  cache:clear
  cache:get
  cache:rebuild
  cache:set
  cache:tags
  cache:warm
```

- First item is `← Back` (styled, always at top)
- Enter on `← Back` returns to level 1
- Enter on a command triggers `onSelect` (loads params form, same as today)
- Esc returns to level 1 (when command list is focused)
- Pane title updates to show path: `"Commands > cache"`

**Implementation approach:**

- Add state to `CommandList`: `level` (1 or 2), `currentNamespace` (string), and store the full `[]NamespaceGroup` data
- Add `showNamespaces()` method that clears the list and populates with namespace items
- Add `showCommands(namespace)` method that clears the list and populates with `← Back` + commands from that group
- Override `handleSelected` to branch on `level`:
  - Level 1: call `showCommands(selectedNamespace)`
  - Level 2, index 0 (Back): call `showNamespaces()`
  - Level 2, index > 0: call existing `onSelect(cmd)`
- `handleChanged` needs updating: at level 1 there are no headers to skip; at level 2, skip index 0 only if navigating past it naturally (or let `← Back` be selectable)

**Focus management changes in `app.go`:**

- Esc when command list is focused AND at level 2: go back to level 1 (not a global Esc handler change — the command list handles it internally)
- Esc when command list is focused AND at level 1: no-op (already at top)
- All other Esc behavior unchanged (form → command list)

**Design decisions:**

- Commands without a colon (like `browse`, `recipe`, `version`) appear under the `"other"` namespace — this already works in `discovery.go`
- `← Back` is a selectable item (not a header), styled with `[yellow]` to stand out
- `q` quits from either level

**Acceptance criteria:**
- [ ] TUI launches showing namespace list with command counts (not the full flat command list)
- [ ] Enter on `cache` shows `← Back` + cache commands
- [ ] Enter on a command loads the params form (Sprint 2 behavior preserved)
- [ ] `← Back` returns to namespace list
- [ ] Esc from level 2 (when command list focused) returns to namespace list
- [ ] Pane title shows `"Commands"` at level 1, `"Commands > cache"` at level 2
- [ ] All existing keyboard shortcuts still work (Tab to form, Esc from form, `q` to quit)
- [ ] `go build ./...` succeeds

---

### S4-2 · Output pane focus + scrolling · [ ]

**Owner:** Cursor Composer (alt: Codex)

**Scope:** `internal/tui/app.go`, `internal/tui/output.go`

The output pane is scrollable in code (`SetScrollable(true)`) but there's no way for the user to focus it or scroll with keys. Fix that.

**Focus cycle — add output pane as the third Tab stop:**

```
Tab:  command list → params form → output pane → command list
Esc:  (any pane) → command list
```

**Scrolling interface (when output pane is focused):**

| Key | Action |
|-----|--------|
| Arrow Up / Arrow Down | Scroll one line up/down |
| Space / Page Down | Scroll one page down (like `more`) |
| b / Page Up | Scroll one page up (like `less`) |
| Home / g | Jump to top of output |
| End / G | Jump to bottom of output |
| Esc | Return focus to command list |

This mirrors the `less` pager that terminal users already know.

**Visual focus indicator:**

- When the output pane is focused, change the border color (e.g., yellow or green) so the user can see which pane has focus
- Reset border color when focus leaves

**Constraints:**

- `q` from output pane should NOT quit (only from command list at level 1)
- The pane should still auto-scroll to top when new output arrives (`ScrollToBeginning` — already implemented)

**Acceptance criteria:**
- [ ] Tab from params form moves focus to output pane
- [ ] Tab from output pane moves focus to command list
- [ ] Arrow Up/Down scrolls output one line at a time
- [ ] Space scrolls one page down; `b` scrolls one page up
- [ ] Home/`g` jumps to top; End/`G` jumps to bottom
- [ ] Esc from output pane returns to command list
- [ ] Output pane border color changes when focused
- [ ] `q` does NOT quit when output pane is focused
- [ ] Long output (e.g., `pm:list`) is scrollable through the full content

---

### S4-3 · Command list search/filter · [ ]

**Owner:** Codex (alt: Sonnet)

**Scope:** `internal/tui/commands.go`, `internal/tui/app.go`

Add `/` key to open a filter that narrows the list as the user types.

**Behavior depends on hierarchy level:**

- **Level 1 (namespaces):** `/` filters namespace names. Typing `ca` shows only `cache`, `checklistapi`. Matches are case-insensitive.
- **Level 2 (commands):** `/` filters commands within the current namespace. Typing `reb` shows only `cache:rebuild`.
- **Esc** clears the filter and returns to the full list (without changing level)
- **Enter** selects the highlighted result (drill into namespace, or open command)
- Filter input appears at the top of the command pane (inline, not a modal)

**Implementation:**

- Add a `tview.InputField` filter bar at the top of the command pane (hidden by default)
- `/` shows the filter bar and focuses it
- On each keystroke (`SetChangedFunc`), re-filter the list
- Esc hides the filter bar, clears the filter text, restores the full list
- Enter hides the filter bar and triggers select on the current item
- Store the unfiltered items so they can be restored on Esc

**Acceptance criteria:**
- [ ] `/` opens filter input at top of command pane
- [ ] Typing narrows the namespace list (level 1) or command list (level 2)
- [ ] Filtering is case-insensitive
- [ ] Esc clears filter and restores full list
- [ ] Enter selects the highlighted item
- [ ] Filter bar is hidden when not in use
- [ ] `go build ./...` succeeds

---

### S4-4 · Required argument validation · [ ]

**Owner:** Codex (alt: Haiku)

**Scope:** `internal/tui/params.go`

Prevent the Run button from executing when required arguments are empty.

- Before calling `onRun`, check that all fields with `* ` prefix labels have non-empty values
- If validation fails, show a red error message below the form header: `"Required: <field name>"`
- Don't clear the form or lose the user's input — just show the error
- Clear the validation error when the user starts typing in the offending field

**Acceptance criteria:**
- [ ] Running `user:create` with empty `name` field shows validation error
- [ ] Validation error clears when user types in the required field
- [ ] Running `user:create` with `name` filled works normally
- [ ] Commands with no required args (like `cache:rebuild`) run without validation errors

---

### S4-5 · Sprint 4 demo script · [ ]

**Owner:** whoever finishes last (or the human)

**Scope:** `scripts/demo-s4.sh`

**WHAT TO LOOK FOR:**
1. TUI launches showing namespace list (not flat command list)
2. Enter on a namespace drills into its commands
3. `← Back` and Esc return to namespace list
4. Select a command, fill params, Run — output appears (Sprint 3 still works)
5. Tab cycles through all three panes (commands → params → output), border highlights focused pane
6. Arrow keys scroll output one line at a time; Space pages down; `b` pages up
7. `/` opens filter, typing narrows the list, Esc clears
8. Run `user:create` with empty name — validation error appears
9. `q` exits from command list

**Acceptance criteria:**
- [ ] `./scripts/demo-s4.sh` runs the program for human review
- [ ] Script documents what to look for

---

## Decisions Made This Sprint

*(Fill in as sprint progresses)*

---

## Acceptance

**Status:** Accepted — ready for user testing
**Date:** 2026-03-23
**Reviewed by:** Luke McCormick

---

## Deferred to Sprint 5

- Error handling (DDEV not running, drush unavailable mid-session, command timeouts) — Phase 9
- Cross-platform build (darwin/arm64, darwin/amd64, linux/amd64) — Phase 10
- GitHub release workflow — Phase 10
- README.md — Phase 10
- Binary download in DDEV add-on shell script

---

## References

- `AGENT.md` — project constraints, UI layout, development process
- `DOC/implementation-plan.md` — Phase 9 and 10 details
- `SPRINTS/sprint-3.md` — previous sprint, see "Deferred" section
- `HANDOFF/handoff-2026-03-20-sprint3-claude.md` — last session handoff
