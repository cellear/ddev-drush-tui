package drush

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

type NamespaceGroup struct {
	Namespace string
	Commands  []Command
}

type Command struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Aliases     []string   `json:"aliases"`
	Definition  Definition `json:"definition"`
}

type Definition struct {
	RawArguments json.RawMessage `json:"arguments"`
}

type commandListEnvelope struct {
	Commands []Command `json:"commands"`
}

var commandCache []NamespaceGroup

func ListCommands() ([]NamespaceGroup, error) {
	if commandCache != nil {
		return commandCache, nil
	}

	cmd := exec.Command("ddev", "drush", "list", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ddev drush list failed: %w", err)
	}

	var envelope commandListEnvelope
	if err := json.Unmarshal(output, &envelope); err != nil {
		return nil, fmt.Errorf("parse drush list: %w", err)
	}

	groupsMap := make(map[string][]Command)
	for _, c := range envelope.Commands {
		if shouldFilter(c.Name) {
			continue
		}

		namespace := "other"
		if strings.Contains(c.Name, ":") {
			namespace = strings.Split(c.Name, ":")[0]
		}

		groupsMap[namespace] = append(groupsMap[namespace], c)
	}

	var groups []NamespaceGroup
	for ns, cmds := range groupsMap {
		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].Name < cmds[j].Name
		})
		groups = append(groups, NamespaceGroup{
			Namespace: ns,
			Commands:  cmds,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Namespace < groups[j].Namespace
	})

	commandCache = groups
	return groups, nil
}

func shouldFilter(name string) bool {
	if strings.HasPrefix(name, "_") {
		return true
	}
	if name == "help" || name == "list" {
		return true
	}
	return false
}
