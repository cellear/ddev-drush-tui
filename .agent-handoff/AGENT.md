# Agent Handoff Protocol

Follow this protocol to maintain context across sessions.

## First-Time Setup

If this project doesn't already have a tool-specific instruction file, ask the user which AI tools they use and create the matching files:

| Tool | Create file |
|------|------------|
| Claude Code | `CLAUDE.md` |
| OpenCode / Codex | `AGENTS.md` |
| Cursor | `.cursorrules` |
| GitHub Copilot | `.github/copilot-instructions.md` |
| Gemini CLI | `.gemini/styleguide.md` |

Each file should contain: `Read and follow AGENT.md in this project's root directory.`

## Directories

- `HANDOFF/` — Session journals, chronological
- `DOC/` — Reference docs, persistent knowledge by topic

Create these directories if they don't exist.

## Starting a Session

1. Read recent files in `HANDOFF/` (newest first)
2. Read relevant files in `DOC/`
3. Summarize current project state for the user
4. Ask what to work on next

## Ending a Session

Create a handoff document before the session ends.

**File:** `HANDOFF/handoff-[yyyy-mm-dd]-[task]-[author].md`

No spaces in filenames. Examples:
- `handoff-2026-01-20-auth-bug-claude.md`
- `handoff-2026-01-21-api-refactor-gemini.md`
- `handoff-2026-01-22-review-jane.md`

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

## Version History

Source: https://github.com/cellear/agent-handoff

- **1.1** (2026-02-16) — Split into README + AGENT.md; added tool-specific setup; simplified DOC guidance
- **1.0** (2026-01-25) — Initial protocol: HANDOFF/ and DOC/ directories, session workflow, naming convention
