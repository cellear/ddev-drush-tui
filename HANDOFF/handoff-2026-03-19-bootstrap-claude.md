# Handoff: Bootstrap — 2026-03-19

## Models Used This Session

| Task | Model |
|------|-------|
| Project planning, architecture, AGENT.md rewrite | Claude Opus |
| DOC files, sprint plan, handoff | Claude Sonnet |

---

## What Was Attempted and Outcome

**Goal:** Bootstrap the ddev-drush-tui project from a copy of the ddev-xdebug-tui repo.

**Outcome: Sprint 0 complete.** All bootstrap tasks done. Ready for Sprint 1.

---

## What Worked

- Moved all xdebug Go source, binaries, and DDEV files to `XDEBUG-TRASHED/` cleanly
- Rewrote `AGENT.md` for the drush-tui project
- Created `.gemini/styleguide.md` for Gemini CLI support
- Updated `.gitignore` (drupal-cms/, cms-2.0.4.zip, XDEBUG-*, .claude/)
- Initialized git repo, made clean initial commit (11 files, zero xdebug artifacts)
- Ran `ddev drush list --format=json` and `ddev drush help <cmd> --format=json` against the drupal-cms test site to capture real JSON output
- Wrote `DOC/drush-discovery.md` documenting JSON structures and the inconsistent `arguments` field issue
- Wrote `DOC/implementation-plan.md` with phased breakdown and owner assignments
- Wrote `SPRINTS/sprint-1.md` with 5 stories, acceptance criteria, and demo checkpoint

---

## What Did Not Work / Was Not Done

- `drupal-cms/` is gitignored (local test harness only) — future agents need it running before testing
- No Go code exists yet — Sprint 0 was documentation and setup only
- Binary download logic in `commands/host/drush-tui` is deferred to Sprint 4 (release prep)

---

## Current State

**Repository:** `/Users/lukemccormick/Sites/DDEV/ddev-drush-tli/`
**Git branch:** `main`, 1 commit (initial bootstrap)
**Remote:** Not yet set. User will push to `github.com/cellear/ddev-drush-tui` when token is refreshed.

**Directory structure:**
```
.agent-handoff/      # Handoff protocol (Claude, Codex, Cursor)
.gemini/             # Gemini CLI instructions
.gitignore           # Updated for drush-tui
AGENT.md             # ← Rewritten for this project
DOC/
  drush-discovery.md # JSON format docs with real examples
  implementation-plan.md # Architecture and phase breakdown
HANDOFF/
  handoff-2026-03-19-bootstrap-claude.md  # (this file)
LEARNINGS/           # Empty, ready for Sprint 1 learnings
SPRINTS/
  sprint-1.md        # Ready to execute
XDEBUG-SAVED/        # Archive of old project docs (gitignored)
XDEBUG-TRASHED/      # Old xdebug source code (gitignored)
drupal-cms/          # Drupal CMS test site (gitignored, local)
```

No Go code yet. Next session starts at S1-1.

---

## Open Questions

1. **GitHub remote:** User needs to refresh the gh token before pushing. Repo will be `github.com/cellear/ddev-drush-tui`.

2. **drupal-cms project name in DDEV:** The test site may have a DDEV project name different from "drupal-cms". Confirm with `ddev describe --json` from that directory. The DDEV context detection code (S1-2) will need to handle whatever name is actually configured.

3. **tview version:** Use the same version as ddev-xdebug-tui. Check `XDEBUG-TRASHED/go.mod` for the exact version, then pin it in the new `go.mod`.

4. **Commands to filter from list:** The filter list in `DOC/drush-discovery.md` is a first pass. Sprint 1 implementation will reveal whether the list is complete. Commands starting with `_` and `help`/`list` are filtered; add more if the command list looks noisy.

---

## Files Created or Modified

- `AGENT.md` — rewritten
- `.gemini/styleguide.md` — created
- `.gitignore` — updated
- `HANDOFF/.gitkeep`, `DOC/.gitkeep`, `LEARNINGS/.gitkeep`, `SPRINTS/.gitkeep` — created (to track empty dirs)
- `DOC/drush-discovery.md` — created
- `DOC/implementation-plan.md` — created
- `SPRINTS/sprint-1.md` — created
- `HANDOFF/handoff-2026-03-19-bootstrap-claude.md` — this file

---

## References

- `DOC/implementation-plan.md` — full architecture, use for all Sprint 1+ work
- `DOC/drush-discovery.md` — read this before implementing discovery.go or help.go
- `XDEBUG-SAVED/DOC/ddev-addon-conventions.md` — critical before implementing S1-5 add-on stub
- `XDEBUG-SAVED/SPRINTS/sprint-1.md` — reference for sprint format (from old project)

---

## Possible Next Steps

Next session: **Sprint 1** — Go scaffold + command discovery.

Suggested order (matching story assignments):
1. **S1-1** (Codex): `go mod init`, add tview, stub all files, `go build ./...` passes
2. **S1-2** (Codex): DDEV context detection, wire to main.go
3. **S1-3** (Gemini CLI): Drush command discovery — read `DOC/drush-discovery.md` first
4. **S1-4** (Cursor Composer): Minimal TUI, wire everything together
5. **S1-5** (Codex + Human): DDEV add-on stub, `make install`, `ddev drush-tui` works

Do not proceed from S1-4 to S1-5 until the Sprint 1 demo checkpoint passes (TUI shows commands, nav works, q exits).
