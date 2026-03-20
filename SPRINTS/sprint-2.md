# Sprint 2: Parameter Discovery + Form

**Sprint Goal:** Selecting a command in the TUI populates the right pane with a dynamic parameter form showing the command's arguments and options. Required arguments are validated. The form is navigable with Tab/arrow keys. No command execution yet — that's Sprint 3.

**Test site:** `drupal-cms/` in the repo root (gitignored; must be running via `ddev start`)

---

## Stories

### S2-1 · Parameter discovery (help parsing) · [ ]

**Owner:** Codex (alt: Sonnet)

**Scope:** `internal/drush/help.go`

**CRITICAL:** Read `DOC/drush-discovery.md` before implementing. The help JSON has `is_required`, `accept_value`, etc. as **strings** ("0" or "1"), not booleans. The `arguments` and `options` fields are at the top level (not nested under `definition` like in `list`).

- Replace the stub `CommandHelp` struct with the full struct from `DOC/drush-discovery.md`
- Run `ddev drush help <command> --format=json`, capture stdout
- Unmarshal into `CommandHelp` struct
- Filter out global Drush options (see filter list in `DOC/drush-discovery.md`)
- Cache result per command name for the session (same pattern as `discovery.go`)
- Export `Help(command string) (*CommandHelp, error)`

Structs needed (from `DOC/drush-discovery.md`):

```go
type CommandHelp struct {
    ID          string              `json:"id"`
    Name        string              `json:"name"`
    Description string              `json:"description"`
    Help        string              `json:"help"`
    Aliases     []string            `json:"aliases"`
    Arguments   map[string]Argument `json:"arguments"`
    Options     map[string]Option   `json:"options"`
    Examples    []Example           `json:"examples"`
}

type Argument struct {
    Name        string `json:"name"`
    IsRequired  string `json:"is_required"`   // "0" or "1"
    IsArray     string `json:"is_array"`      // "0" or "1"
    Description string `json:"description"`
}

type Option struct {
    Name            string   `json:"name"`
    Shortcut        string   `json:"shortcut"`
    AcceptValue     string   `json:"accept_value"`      // "0" or "1"
    IsValueRequired string   `json:"is_value_required"`  // "0" or "1"
    IsMultiple      string   `json:"is_multiple"`        // "0" or "1"
    Description     string   `json:"description"`
    Defaults        []string `json:"defaults"`
}

type Example struct {
    Usage       string `json:"usage"`
    Description string `json:"description"`
}
```

**Global options to filter out:**

```
druplicon, notify, xh-link, help, silent, quiet, verbose, debug,
yes, no, uri, root, simulate, backend, pipe, interact, bootstrap,
config, alias-path, include, phpunit-path, generate-md, symfony
```

**Acceptance criteria:**
- [ ] `drush.Help("cache:rebuild")` returns populated `CommandHelp` with filtered options
- [ ] `drush.Help("user:create")` returns `CommandHelp` with the `name` argument (required)
- [ ] Global options are not present in the returned `Options` map
- [ ] Results are cached — calling `Help()` twice for the same command runs the subprocess only once
- [ ] `go build ./...` still succeeds

---

### S2-2 · Wire command selection to help lookup · [ ]

**Owner:** Codex (alt: Sonnet)

**Scope:** `internal/tui/app.go`

This is a small wiring story. When the user selects a command (Enter) in the command list, call `drush.Help()` and pass the result to the params pane.

- Update the `onSelect` callback in `NewApp` to call `drush.Help(cmd.Name)`
- On success, call a method on the params pane to display the result (S2-3 implements the display; this story just passes the data)
- On error, show the error in the params pane text
- Show the command description and aliases in the params pane header area

**Acceptance criteria:**
- [ ] Selecting a command triggers `drush.Help()` (visible via debug output or params pane update)
- [ ] Errors from `Help()` are shown in the params pane, not swallowed silently
- [ ] `go build ./...` still succeeds

---

### S2-3 · Parameter form rendering · [ ]

**Owner:** Cursor Composer (alt: Haiku)

**Scope:** `internal/tui/params.go`

**Read:** `AGENT.md` → User Interface section for the layout spec. `DOC/drush-discovery.md` → option/argument field meanings.

Replace the `ParamsView` placeholder with a dynamic form:

- Change `ParamsView` to wrap a `tview.Form` (instead of `tview.TextView`)
- Add a method `ShowParams(help *drush.CommandHelp)` that clears and rebuilds the form:
  - **Header area:** command name, description, aliases (as static text above the form)
  - **Required arguments:** `tview.InputField` with `* ` prefix in label. One per argument where `is_required == "1"`.
  - **Optional arguments:** `tview.InputField` without prefix. Arguments where `is_required == "0"`.
  - **Options that accept values** (`accept_value == "1"`): `tview.InputField`. Pre-fill with default if `defaults` is present.
  - **Boolean flags** (`accept_value == "0"`): `tview.Checkbox`.
  - **"Run" button** at bottom (callback is no-op for now — Sprint 3 wires to executor)
  - **"Cancel" button** clears form and returns focus to command list
- Add a method `ShowError(err error)` that displays error text
- Add a method `ShowPlaceholder()` that shows "Select a command to see its parameters."
- Arguments should appear before options in the form

**Focus management:**
- Tab cycles between command list and params form
- When params form is focused, Tab/Enter navigate within the form fields
- Esc returns focus to command list
- `q` should NOT quit when the form is focused (only when command list is focused — already handled)

**Acceptance criteria:**
- [ ] Selecting `user:create` shows a form with `* name` (required input field)
- [ ] Selecting `cache:rebuild` shows a form with the `cache-clear` option (pre-filled with default)
- [ ] Boolean flags render as checkboxes
- [ ] Global options (druplicon, verbose, etc.) are not shown
- [ ] "Cancel" returns focus to command list
- [ ] Tab moves focus between command list and params form
- [ ] Esc from params form returns focus to command list
- [ ] Form clears and rebuilds when a different command is selected

---

### S2-4 · Sprint 2 demo script · [ ]

**Owner:** whoever finishes S2-3 (or the human)

**Scope:** `scripts/demo-s2.sh`

Create a demo script following the pattern in `scripts/demo-s1.sh`.

**WHAT TO LOOK FOR:**
1. TUI launches with grouped command list (Sprint 1 still works)
2. Selecting `user:create` shows a form with `* name` required field
3. Selecting `cache:rebuild` shows options without global Drush options
4. Tab cycles focus between command list and form
5. Esc returns from form to command list
6. Cancel button works
7. `q` exits from command list

**Acceptance criteria:**
- [ ] `./scripts/demo-s2.sh` runs the program for human review
- [ ] Script documents what to look for

---

## Decisions Made This Sprint

*(Fill in as sprint progresses)*

---

## Deferred to Sprint 3

- Command execution (`drush/executor.go`)
- Output pane (`tui/output.go`)
- "Run" button functionality (button exists but is no-op)
- Error handling for DDEV/Drush failures mid-session

---

## References

- `AGENT.md` — project constraints, UI layout, development process
- `DOC/implementation-plan.md` — Phase 6 and 7 details
- `DOC/drush-discovery.md` — JSON format docs, struct designs, global option filter list
- `SPRINTS/sprint-1.md` — completed sprint for format reference
