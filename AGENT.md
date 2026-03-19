
# AGENT.md

Guidelines for AI agents working on **ddev-drush-tui**

---

# Project Overview

This project implements a **terminal UI for the Drush command-line tool**, running inside DDEV.

The TUI lets users:

- browse available Drush commands grouped by namespace
- select a command and see its arguments and options
- fill in parameters via a guided form
- execute the command and see output

Users don't need to memorize Drush syntax. DDEV context (project, container) is handled automatically.

---

# Core Design Philosophy

This project prioritizes:

1. **Simplicity**
2. **Readability**
3. **Minimal dependencies**
4. **Small codebase**
5. **Predictable behavior**

The maintainer intends to read and understand every line of code before release.

Agents must **avoid introducing complexity** unless explicitly requested.

---

# Absolute Constraints

Do NOT introduce:

- command auto-execution without user confirmation
- persistent configuration or history (PoC scope)
- multi-site handling (one DDEV project at a time)
- concurrency unless absolutely necessary
- heavy frameworks
- complex architecture patterns

Keep the interaction model simple: browse, select, fill, run, read output.

---

# Technology Stack

Language: Go
UI Framework: tview
Drush Interaction: subprocess calls via `ddev drush`
Distribution: DDEV add-on invoking a Go binary

---

# How the TUI Interacts with Drush

All Drush interaction goes through `ddev drush`. The TUI runs on the host, not in the container.

**Command discovery:**
`ddev drush list --format=json` returns all available commands with descriptions, grouped by namespace.

**Parameter discovery:**
`ddev drush help <command> --format=json` returns arguments, options, aliases, and descriptions for a specific command.

**Command execution:**
`ddev drush <command> [arguments] [--options]` runs the command. Stdout and stderr are captured and displayed in the output pane.

Discovery is done at startup (command list) and on-demand (parameters when a command is selected). Results are cached for the session.

---

# Architecture

```
cmd/
  ddev-drush-tui/
    main.go                 # Entry point

internal/
  ddev/
    context.go              # Detect DDEV project, verify drush is available
  drush/
    discovery.go            # Parse `ddev drush list --format=json`
    help.go                 # Parse `ddev drush help <cmd> --format=json`
    executor.go             # Run commands, capture output
  tui/
    app.go                  # tview application, layout, keybindings
    commands.go             # Command list pane (left)
    params.go               # Parameter form pane (right)
    output.go               # Command output pane (bottom)
```

Responsibilities:

**ddev** — detect whether we're in a DDEV project, get project name, verify drush is available

**drush/discovery** — run `ddev drush list`, parse JSON into Go structs, group by namespace

**drush/help** — run `ddev drush help <cmd>`, parse JSON into argument/option structs

**drush/executor** — build command string from user input, run via `ddev drush`, capture output

**tui** — render the three-pane layout, handle navigation, wire up form submission to executor

---

# User Interface

The UI is a **three-pane terminal layout**:

```
+------------------------------------------------------------+
| ddev-drush-tui | project: drupal-cms                       |
+--------------------+---------------------------------------+
| Commands           | Parameters                            |
|                    |                                       |
| > cache:rebuild    |  name (required): ___                 |
|   config:export    |  --verbose  [yes/no]                  |
|   config:import    |                                       |
|   core:status      |  [ Run ]    [ Cancel ]                |
|   pm:list          |                                       |
|   ...              |                                       |
+--------------------+---------------------------------------+
| Output                                                     |
|                                                            |
| $ ddev drush cache:rebuild                                 |
| [success] Cache rebuild complete.                          |
|                                                            |
+------------------------------------------------------------+
```

Left pane: scrollable list of commands, grouped by namespace
Right pane: dynamic form showing arguments and options for selected command
Bottom pane: scrollable output from the last executed command

---

# Key Bindings

Tab     cycle focus between panes
Enter   select command / submit form
Esc     cancel parameter form, return to command list
q       quit (when command list is focused)
/       filter commands (search)

---

# Concurrency Policy

Prefer **single-threaded execution**.

If goroutines are used, they must be minimal and clearly documented.

The main use case for a goroutine: running a Drush command without freezing the UI.

---

# Code Quality Requirements

All code must be:

- small
- readable
- well-commented
- free of unnecessary abstraction

Prefer explicit code over clever code.

---

# Implementation Order

Agents should implement features in this order:

1. DDEV context detection
2. Drush command discovery (JSON parsing)
3. Minimal TUI with command list
4. Parameter discovery (help parsing)
5. Parameter form rendering
6. Command execution
7. Output display
8. DDEV add-on packaging

Each step should compile and run before moving to the next.

---

# Testing Strategy

Manual testing is acceptable for the PoC.

Typical test loop:

1. Launch TUI from a DDEV project directory
2. Browse commands
3. Select a command
4. Fill in parameters
5. Run the command
6. Verify output matches `ddev drush <cmd>` run directly

The `drupal-cms/` directory in this repo is the test site.

---

# Things To Avoid

Agents should not:

- introduce large dependencies
- introduce unnecessary abstractions
- rewrite the architecture
- add new features not requested
- attempt to handle commands that require interactive input (those are out of PoC scope)

If unsure, prefer the simplest possible implementation.

---

# Success Criteria

The project is successful if a developer can:

1. run `ddev drush-tui`
2. see a list of available Drush commands
3. select a command
4. fill in its parameters
5. run it
6. see the output

without knowing Drush syntax.

---

# Future Features (Not PoC)

These may be implemented later:

- command favorites / recently used
- persistent configuration
- multi-site support
- command history across sessions
- interactive input passthrough

Agents should NOT attempt these during PoC development.

---

# Development Process

This project uses the **agent handoff protocol**. See `.agent-handoff/AGENT.md` for the full protocol.

Key directories:

- `HANDOFF/` — session journals, chronological
- `DOC/` — persistent reference docs
- `SPRINTS/` — sprint plans with user stories
- `LEARNINGS/` — end-of-sprint knowledge capture

Agents must read recent HANDOFF/ files before starting work.
Agents must write a handoff document before ending a session.

Sprint plans live in `SPRINTS/sprint-N.md` with user stories, acceptance criteria, and `[done]` tags.

---

# Tool-Specific Instructions

| Tool | File |
|------|------|
| Claude Code | `.agent-handoff/CLAUDE.md` |
| Codex | `.agent-handoff/AGENTS.md` |
| Cursor | `.agent-handoff/.cursorrules` |
| Gemini CLI | `.gemini/styleguide.md` |

All point to this file (AGENT.md) as the source of truth.
