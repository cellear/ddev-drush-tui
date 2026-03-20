
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

## Key Directories

- `HANDOFF/` — session journals, chronological
- `DOC/` — persistent reference docs
- `SPRINTS/` — sprint plans with user stories
- `LEARNINGS/` — end-of-sprint knowledge capture

Sprint plans live in `SPRINTS/sprint-N.md` with user stories, acceptance criteria, and `[done]` tags.

## Commit Attribution

Every commit made by an AI assistant MUST include a `Co-Authored-By` trailer identifying the model. Examples:

```
Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
Co-Authored-By: Codex <noreply@openai.com>
Co-Authored-By: Gemini CLI <noreply@google.com>
Co-Authored-By: Cursor Composer <noreply@cursor.com>
```

## Starting a Session

1. Read this file (AGENT.md) completely
2. Read recent files in `HANDOFF/` (newest first)
3. Read relevant files in `DOC/`
4. Summarize current project state for the user
5. Ask what to work on next

## Ending a Session (Handoff Protocol)

Create a handoff document before the session ends.

**File:** `HANDOFF/handoff-[yyyy-mm-dd]-[task]-[author].md`

No spaces in filenames. Examples:
- `handoff-2026-01-20-auth-bug-claude.md`
- `handoff-2026-01-21-api-refactor-gemini.md`

**Include:**
- What was attempted and the outcome
- What worked, what didn't
- Current state and blockers
- Open questions
- Files created or modified
- References to `DOC/` files and prior handoffs
- **Next-assistant prompt** (see below)

## Next-Assistant Prompt

Every handoff MUST end with a section called `## Prompt for Next Assistant`. This is a ready-to-paste prompt the human can give to the next AI assistant to start their session. It should:

1. State which sprint story to work on next (e.g., "You are working on S1-4")
2. List the files to read first (AGENT.md, the sprint plan, relevant DOC/ files, the previous handoff)
3. Summarize what's already done and what's expected
4. Remind the assistant of key constraints (e.g., "read DOC/drush-discovery.md before implementing")
5. Remind the assistant to sign commits with `Co-Authored-By`

After writing the handoff file, **tell the human directly:**
- The exact prompt to paste (in a code block)
- Which assistant is scheduled to receive it (from the sprint plan's Owner field)
- What story it covers

## Updating DOC Files

When you learn something persistent about the project — architecture decisions, deploy steps, conventions, known issues — update or create a relevant file in `DOC/`. Prefer updating existing files over creating new ones.

When updating a DOC file, add or update a `Last updated: yyyy-mm-dd by [author]` line at the bottom.

---

# Tool-Specific Instructions

Each tool has a pointer file in the repo root that says "Read and follow AGENT.md":

| Tool | File |
|------|------|
| Claude Code | `CLAUDE.md` |
| Codex | `AGENTS.md` |
| Cursor | `.cursorrules` |
| Gemini CLI | `.gemini/styleguide.md` |
