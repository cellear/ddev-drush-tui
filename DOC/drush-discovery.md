# Drush Discovery: JSON Formats

This document captures the real JSON output from `ddev drush list` and `ddev drush help <command>`, documenting the structures our Go code will parse.

Captured from: `drupal-cms/` test site (Drupal CMS 2.0.4, Drush 13.7.1)

---

## Command List: `ddev drush list --format=json`

The full output is an array of command objects at the top level (not namespaced — all flat).

### Top-level structure

```json
{
  "application": {
    "name": "Drush Commandline Tool",
    "version": "13.7.1.0"
  },
  "commands": [
    { ... command object ... },
    { ... command object ... }
  ]
}
```

### Command object (minimal, no arguments)

```json
{
  "name": "cache:rebuild",
  "description": "Rebuild all caches.",
  "usage": ["cache:rebuild [--cache-clear [CACHE-CLEAR]] ..."],
  "help": "This is a copy of core/rebuild.php.",
  "aliases": ["cr", "rebuild", "cache-rebuild"],
  "definition": {
    "arguments": [],
    "options": { ... }
  }
}
```

### Command object (with required argument)

```json
{
  "name": "user:create",
  "description": "Create a user account.",
  "usage": ["user:create [--format [FORMAT]] ... [--] <name>"],
  "definition": {
    "arguments": {
      "name": {
        "name": "name",
        "is_required": true,
        "is_array": false,
        "description": "The name of the account to add",
        "default": null
      }
    },
    "options": { ... }
  }
}
```

### Key fields for command discovery

| Field | Type | Notes |
|-------|------|-------|
| `name` | string | Full name, e.g. `cache:rebuild`. Namespace = part before `:` |
| `description` | string | Short description for command list display |
| `aliases` | []string | Short aliases, e.g. `cr` for `cache:rebuild` |
| `definition.arguments` | map or [] | Empty array `[]` when no arguments; object map when there are args |
| `definition.options` | map | Always present; includes Drush global options (druplicon, notify, etc.) |

### Namespace extraction

The namespace is the part before the colon in the command name:

- `cache:rebuild` → namespace `cache`
- `user:create` → namespace `user`
- `core:status` → namespace `core`
- `_complete` → no colon → treat as namespace `system` or group as "other"

---

## Command Help: `ddev drush help <command> --format=json`

This returns richer detail for a specific command. Used when the user selects a command to see its parameters.

### Structure

```json
{
  "id": "cache:rebuild",
  "name": "cache:rebuild",
  "usages": ["cache:rebuild [--cache-clear [CACHE-CLEAR]] ..."],
  "description": "Rebuild all caches.",
  "help": "This is a copy of core/rebuild.php.",
  "aliases": ["cr", "rebuild", "cache-rebuild"],
  "arguments": {},
  "options": {
    "cache-clear": { ... option object ... }
  },
  "examples": [
    {
      "usage": "drush cache:rebuild",
      "description": "Rebuild all caches."
    }
  ]
}
```

Note: unlike `list`, the help output has `arguments` and `options` at the top level (not nested under `definition`).

### Option object

```json
{
  "name": "--cache-clear",
  "shortcut": "",
  "accept_value": "1",
  "is_value_required": "0",
  "is_multiple": "0",
  "description": "Set to 0 to suppress normal cache clearing...",
  "defaults": ["true"]
}
```

### Argument object (from `user:create`)

```json
{
  "name": "name",
  "is_required": "1",
  "is_array": "0",
  "description": "The name of the account to add"
}
```

Note: `is_required`, `is_array`, `is_multiple`, `accept_value` are **strings** ("0" or "1"), not booleans. Parse accordingly.

### Key fields for parameter form

| Field | Type | Notes |
|-------|------|-------|
| `arguments.<key>.name` | string | Positional argument name |
| `arguments.<key>.is_required` | string "0"/"1" | Whether the argument is required |
| `arguments.<key>.description` | string | Help text to show in form |
| `options.<key>.name` | string | Option name including `--` prefix |
| `options.<key>.accept_value` | string "0"/"1" | "0" = boolean flag, "1" = takes a value |
| `options.<key>.is_value_required` | string "0"/"1" | "1" = value is mandatory if option is used |
| `options.<key>.description` | string | Help text |
| `options.<key>.defaults` | []string or absent | Default values |
| `examples` | []object | Usage examples to show in the form |

---

## Global Options to Filter Out

Every Drush command includes these global options that are not useful to show in the TUI parameter form:

- `--druplicon` — ASCII art easter egg
- `--notify` — desktop notification
- `--xh-link` — XHProf profiling
- `--help`, `--silent`, `--quiet`, `--verbose`, `--debug`
- `--yes`, `--no`
- `--uri`, `--root`

The TUI should filter these from the parameter form and only show command-specific arguments and options.

A reasonable filter: show only options that appear in `definition.options` in the `list` output but **not** in the base Drush global option set. Or simply: exclude known global option names.

**Known global options to always suppress:**

```
druplicon, notify, xh-link, help, silent, quiet, verbose, debug,
yes, no, uri, root, simulate, backend, pipe, interact, bootstrap,
config, alias-path, include, phpunit-path, generate-md, symfony
```

---

## Struct Design for Go

### For command list (from `ddev drush list`)

```go
type CommandList struct {
    Application AppInfo   `json:"application"`
    Commands    []Command `json:"commands"`
}

type AppInfo struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

type Command struct {
    Name        string     `json:"name"`
    Description string     `json:"description"`
    Aliases     []string   `json:"aliases"`
    Definition  Definition `json:"definition"`
}

type Definition struct {
    Arguments map[string]Argument `json:"arguments"`
    Options   map[string]Option   `json:"options"`
}
```

Note: `definition.arguments` is `[]` (empty array) when there are no arguments, but an object map when there are. This is inconsistent JSON. In Go, use `json.RawMessage` and handle both cases, or use a custom unmarshaler. See the parsing note below.

### For command help (from `ddev drush help`)

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
    IsRequired  string `json:"is_required"`  // "0" or "1"
    IsArray     string `json:"is_array"`     // "0" or "1"
    Description string `json:"description"`
}

type Option struct {
    Name            string   `json:"name"`
    Shortcut        string   `json:"shortcut"`
    AcceptValue     string   `json:"accept_value"`     // "0" or "1"
    IsValueRequired string   `json:"is_value_required"` // "0" or "1"
    IsMultiple      string   `json:"is_multiple"`       // "0" or "1"
    Description     string   `json:"description"`
    Defaults        []string `json:"defaults"`
}

type Example struct {
    Usage       string `json:"usage"`
    Description string `json:"description"`
}
```

---

## Parsing Note: Inconsistent `arguments` Field

In `ddev drush list --format=json`, the `definition.arguments` field is:
- `[]` (empty JSON array) when the command has no arguments
- `{ "name": {...}, ... }` (JSON object) when it has arguments

This is a known Symfony Console quirk. In Go, handle it with `json.RawMessage`:

```go
type Definition struct {
    RawArguments json.RawMessage     `json:"arguments"`
    Options      map[string]Option   `json:"options"`
}

func (d *Definition) Arguments() map[string]Argument {
    var args map[string]Argument
    if err := json.Unmarshal(d.RawArguments, &args); err != nil {
        return nil // was an empty array []
    }
    return args
}
```

The `help` output does NOT have this issue — `arguments` is always an object (empty `{}` when none).

---

*Last updated: 2026-03-19 by claude-sonnet*
