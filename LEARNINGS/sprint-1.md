# Sprint 1 Learnings

- 2026-03-19: On this machine, `ddev describe` uses `--json-output` instead of `--json`. The JSON payload includes a top-level envelope plus a `raw` object; the project name and app root are available under `raw.name` and `raw.approot`.
- 2026-03-19: Drush JSON output (from Symfony Console) has an inconsistent `definition.arguments` field: it's an empty array `[]` when no arguments exist, but an object map `{}` when they do. In Go, this must be handled using `json.RawMessage` to avoid unmarshaling errors.
- 2026-03-19: Internal commands like `_complete`, `help`, and `list` should be filtered out to provide a clean list for the TUI. Commands without a colon in their name are grouped under the "other" namespace by convention.
