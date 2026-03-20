# Sprint 1: Scaffold + Command Discovery

**Sprint Goal:** `ddev drush-tui` (or `go run ./cmd/ddev-drush-tui` from the project directory) launches a TUI showing a live, grouped list of Drush commands fetched from the DDEV test site. Navigation works. `q` exits.

**Test site:** `drupal-cms/` in the repo root (gitignored; must be running via `ddev start`)

---

## Stories

### S1-1 · Go module + directory structure · [x]

**Owner:** Codex (alt: Haiku)

**Scope:**
- `go mod init github.com/cellear/ddev-drush-tui`
- `go get github.com/rivo/tview` and `go get github.com/gdamore/tcell/v2`
- Create directory structure:
  ```
  cmd/ddev-drush-tui/main.go
  internal/ddev/context.go
  internal/drush/discovery.go
  internal/drush/help.go
  internal/drush/executor.go
  internal/tui/app.go
  internal/tui/commands.go
  internal/tui/params.go
  internal/tui/output.go
  ```
- All files are stubs that compile. `main.go` prints "ddev-drush-tui" and exits.

**Acceptance criteria:**
- [x] `go build ./...` succeeds with zero errors
- [x] `go run ./cmd/ddev-drush-tui` prints something and exits cleanly
- [x] `go.mod` shows module `github.com/cellear/ddev-drush-tui`

---

### S1-2 · DDEV context detection · [x]

**Owner:** Codex (alt: Gemini CLI)

**Scope:** `internal/ddev/context.go`

- Define `Context` struct: `{ ProjectName string, AppRoot string }`
- Run `ddev describe --json-output` and parse the project name from the output
- Run `ddev drush version --format=string` to verify drush is available (exit 0)
- If not in a DDEV project: print `"Error: not in a DDEV project directory"` and exit 1
- If drush not available: print `"Error: drush is not available in this DDEV project"` and exit 1
- Export `Detect() (*Context, error)` function

Wire to `main.go`: call `ddev.Detect()`, print project name, then exit.

**Acceptance criteria:**
- [x] Running from `drupal-cms/` directory: prints `"Project: drupal-cms"` (or similar)
- [x] Running from a non-DDEV directory: prints clear error and exits 1
- [x] Running from a DDEV project without drush: prints clear error and exits 1

**Note:** `ddev describe --json-output` output includes a `name` field with the project name. Check the actual JSON by running it against the test site.

---

### S1-3 · Drush command discovery · [x]

**Owner:** Gemini CLI (alt: Sonnet)

**Scope:** `internal/drush/discovery.go`

**CRITICAL:** Read `DOC/drush-discovery.md` before implementing. The `definition.arguments` field in the list output is inconsistently typed (`[]` vs `{}`). Must be handled with `json.RawMessage`.

- Run `ddev drush list --format=json`, capture stdout
- Unmarshal using structs documented in `DOC/drush-discovery.md`
- Extract namespace from command name (everything before `:`); commands without `:` go in namespace `"other"`
- Return `[]NamespaceGroup` where each group has a name and a list of `Command` structs
- Cache result — don't re-run on every keystroke
- Export `ListCommands() ([]NamespaceGroup, error)`

Structs needed:
```go
type NamespaceGroup struct {
    Namespace string
    Commands  []Command
}

type Command struct {
    Name        string
    Description string
    Aliases     []string
}
```

**Acceptance criteria:**
- [x] `discovery.ListCommands()` returns populated groups
- [x] `cache:*` commands appear in the `cache` namespace group
- [x] Commands are sorted alphabetically within each group
- [x] Groups are sorted alphabetically by namespace
- [x] `_complete`, `help`, `list` — internal commands — are filtered out

**Filtering:** Exclude commands whose name starts with `_` and the built-in `help` and `list` commands. These are Symfony Console internals, not Drush commands.

---

### S1-4 · Minimal TUI with command list · [ ]

**Owner:** Cursor Composer (alt: Haiku)

**Scope:** `internal/tui/app.go`, `internal/tui/commands.go`

**Read:** `AGENT.md` → User Interface section for the layout spec.

Layout:
```
+------------------------------------------------------------+
| ddev-drush-tui | project: drupal-cms                       |
+--------------------+---------------------------------------+
| Commands           | (right pane placeholder)               |
| [namespace header] |                                       |
| > cache:rebuild    | Select a command to see parameters    |
|   config:export    |                                       |
+--------------------+---------------------------------------+
| (output placeholder)                                       |
| Output will appear here after running a command.          |
+------------------------------------------------------------+
```

`internal/tui/app.go`:
- `tview.Application` with `tview.Grid` layout (3 rows: header, panels, output)
- Top panels split: `tview.Flex` with commands list (60%) and params placeholder (40%)
- Header `tview.TextView` showing `" ddev-drush-tui | project: <name>"`
- Output `tview.TextView` placeholder
- `q` quits when commands pane is focused

`internal/tui/commands.go`:
- `tview.List` populated with namespace groups
- Namespace names shown as non-selectable, styled differently (e.g., `[yellow]cache[-]`)
- Commands shown as selectable items
- Arrow keys navigate
- Selecting a command: call a registered callback (just log for now; wired to params in Sprint 2)

`internal/tui/params.go` (stub):
- `tview.TextView` showing "Select a command to see its parameters."

`internal/tui/output.go` (stub):
- `tview.TextView` showing "Output will appear here after running a command."

**Wire in `main.go`:** Call `ddev.Detect()`, call `drush.ListCommands()`, call `tui.NewApp(context, commands).Run()`.

**Acceptance criteria:**
- [ ] TUI launches from `drupal-cms/` directory showing Drush commands grouped by namespace
- [ ] Namespace headers are visually distinct (styled, not selectable)
- [ ] Arrow key navigation moves between commands
- [ ] Selecting a command doesn't crash (callback called, no-op for now)
- [ ] `q` exits cleanly
- [ ] Header shows correct project name

**Sprint 1 Demo checkpoint — do not proceed to Sprint 2 until this passes:**
1. `cd drupal-cms && ddev drush-tui` (or `go run ../cmd/ddev-drush-tui`)
2. TUI shows grouped Drush commands
3. Developer scrolls with arrow keys
4. `q` exits

---

### S1-5 · DDEV add-on stub · [ ]

**Owner:** Codex (alt: Haiku) + Human (testing)

**Scope:** `install.yaml`, `commands/host/drush-tui`, `Makefile`

**CRITICAL: Read `XDEBUG-SAVED/DOC/ddev-addon-conventions.md` before implementing.** The old xdebug project hard-won these rules.

Key rules:
- Source files stored WITHOUT `.ddev/` prefix in the repo
- `install.yaml` lists files as `commands/host/drush-tui` (not `.ddev/commands/...`)
- Shell script must have `#ddev-generated` as line 2

`commands/host/drush-tui`:
```bash
#!/bin/bash
## Description: Launch the Drush TUI
## Usage: drush-tui
## Example: ddev drush-tui
#ddev-generated

# Download binary if not present, then run it
# (binary download logic from GitHub releases — placeholder for now)
# For local dev: just run the binary from PATH
ddev-drush-tui
```

`install.yaml`:
```yaml
name: ddev-drush-tui
project_files:
  - commands/host/drush-tui
```

`Makefile`:
```makefile
BINARY=ddev-drush-tui
build:
    go build -o bin/$(BINARY) ./cmd/ddev-drush-tui
install: build
    cp bin/$(BINARY) ~/go/bin/$(BINARY)
dist:
    GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY)-darwin-arm64 ./cmd/ddev-drush-tui
    GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64 ./cmd/ddev-drush-tui
    GOOS=linux  GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64  ./cmd/ddev-drush-tui
clean:
    rm -rf bin/ dist/
```

**Human testing step:**
1. `make install` — builds and copies binary to `~/go/bin/`
2. From `drupal-cms/`: `ddev add-on get ../` (parent directory)
3. `ddev restart`
4. `ddev drush-tui`

**Acceptance criteria:**
- [ ] `make build` succeeds, binary created at `bin/ddev-drush-tui`
- [ ] `make install` copies binary to `~/go/bin/ddev-drush-tui`
- [ ] `ddev add-on get ../` installs without errors (human step)
- [ ] `ddev drush-tui` launches the TUI (human step)

---

## Decisions Made This Sprint

- **Demo scripts, not automated tests:** The `drupal-cms/acceptance-test*.sh` scripts are _demos_, not automated pass/fail tests. They run the program and display output for a human to review. The human decides whether the sprint's work is accepted. Future scripts should follow this pattern: show what to look for, run the program, let the human judge.

---

## Deferred to Later Sprints

- Parameter form (Sprint 2)
- Command execution (Sprint 3)
- Search/filter (Sprint 4)
- Binary download in shell script (Sprint 4 / release)

---

## References

- `DOC/implementation-plan.md` — full architecture and phase breakdown
- `DOC/drush-discovery.md` — JSON format documentation with real examples
- `XDEBUG-SAVED/DOC/ddev-addon-conventions.md` — critical DDEV add-on rules
- `AGENT.md` — project constraints and TUI layout spec
