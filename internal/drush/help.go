package drush

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// CommandHelp holds the parsed output of "ddev drush help <command> --format=json".
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

// globalOptions lists Drush global options to filter from the parameter form.
var globalOptions = map[string]bool{
	"druplicon":    true,
	"notify":       true,
	"xh-link":      true,
	"help":         true,
	"silent":       true,
	"quiet":        true,
	"verbose":      true,
	"debug":        true,
	"yes":          true,
	"no":           true,
	"uri":          true,
	"root":         true,
	"simulate":     true,
	"backend":      true,
	"pipe":         true,
	"interact":     true,
	"bootstrap":    true,
	"config":       true,
	"alias-path":   true,
	"include":      true,
	"phpunit-path": true,
	"generate-md":  true,
	"symfony":      true,
}

// helpCache stores parsed help results keyed by command name.
var helpCache = map[string]*CommandHelp{}

// Help fetches and parses help JSON for a single Drush command.
// Results are cached per command name for the session.
func Help(command string) (*CommandHelp, error) {
	if cached, ok := helpCache[command]; ok {
		return cached, nil
	}

	cmd := exec.Command("ddev", "drush", "help", command, "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ddev drush help %s failed: %w", command, err)
	}

	var help CommandHelp
	if err := json.Unmarshal(output, &help); err != nil {
		return nil, fmt.Errorf("parse drush help %s: %w", command, err)
	}

	// Filter out global Drush options.
	for key := range help.Options {
		if globalOptions[key] {
			delete(help.Options, key)
		}
	}

	helpCache[command] = &help
	return &help, nil
}
