package tui

import (
	"fmt"

	"github.com/cellear/ddev-drush-tui/internal/drush"
	"github.com/rivo/tview"
)

const (
	cmdLevelNamespaces = 1
	cmdLevelCommands   = 2
)

// CommandList is a tview.List with two-level navigation: namespaces (with
// counts) then commands within a namespace. Enter drills in; ← Back, Esc, or
// Enter on Back returns to the namespace list.
type CommandList struct {
	*tview.List
	groups           []drush.NamespaceGroup
	level            int
	currentNamespace string
	onSelect         func(cmd *drush.Command)
}

// NewCommandList builds a namespace-first list. onSelect is called when the
// user chooses a command at level 2; nil is safe.
func NewCommandList(groups []drush.NamespaceGroup, onSelect func(cmd *drush.Command)) *CommandList {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetUseStyleTags(true, false)

	cl := &CommandList{
		List:     list,
		groups:   groups,
		level:    cmdLevelNamespaces,
		onSelect: onSelect,
	}
	cl.showNamespaces()
	list.SetSelectedFunc(cl.handleSelected)
	return cl
}

func (cl *CommandList) showNamespaces() {
	cl.level = cmdLevelNamespaces
	cl.currentNamespace = ""
	cl.Clear()
	cl.SetTitle("Commands")
	for _, g := range cl.groups {
		label := fmt.Sprintf("%s  (%d)", g.Namespace, len(g.Commands))
		cl.AddItem(label, "", 0, nil)
	}
	if cl.GetItemCount() > 0 {
		cl.SetCurrentItem(0)
	}
}

func (cl *CommandList) showCommands(namespace string) {
	cl.level = cmdLevelCommands
	cl.currentNamespace = namespace
	cl.Clear()
	cl.SetTitle("Commands > " + namespace)
	cl.AddItem("[yellow]← Back[-]", "", 0, nil)

	var cmds []drush.Command
	for i := range cl.groups {
		if cl.groups[i].Namespace == namespace {
			cmds = cl.groups[i].Commands
			break
		}
	}
	for i := range cmds {
		c := cmds[i]
		cl.AddItem(c.Name, "", 0, nil)
	}
	if cl.GetItemCount() > 0 {
		cl.SetCurrentItem(0)
	}
}

func (cl *CommandList) handleSelected(index int, mainText, secondaryText string, shortcut rune) {
	if cl.level == cmdLevelNamespaces {
		if index < 0 || index >= len(cl.groups) {
			return
		}
		cl.showCommands(cl.groups[index].Namespace)
		return
	}
	if index == 0 {
		cl.showNamespaces()
		return
	}
	cmd := cl.commandAtLevel2(index)
	if cmd != nil && cl.onSelect != nil {
		cl.onSelect(cmd)
	}
}

func (cl *CommandList) commandAtLevel2(index int) *drush.Command {
	if index <= 0 {
		return nil
	}
	cmdIdx := index - 1
	for i := range cl.groups {
		if cl.groups[i].Namespace != cl.currentNamespace {
			continue
		}
		g := &cl.groups[i]
		if cmdIdx < len(g.Commands) {
			return &g.Commands[cmdIdx]
		}
		break
	}
	return nil
}

// BackToNamespaces moves from level 2 to level 1. It returns true if the list
// was at level 2 and was updated.
func (cl *CommandList) BackToNamespaces() bool {
	if cl.level != cmdLevelCommands {
		return false
	}
	cl.showNamespaces()
	return true
}
