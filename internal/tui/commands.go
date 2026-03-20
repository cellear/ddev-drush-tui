package tui

import (
	"github.com/cellear/ddev-drush-tui/internal/drush"
	"github.com/rivo/tview"
)

// CommandList is a tview.List showing Drush commands grouped by namespace.
// Namespace headers are styled and non-selectable (arrow keys skip over them).
type CommandList struct {
	*tview.List
	headerIndices map[int]bool      // indices that are namespace headers
	commandsByIdx []*drush.Command  // command at each index (nil for headers)
	onSelect      func(cmd *drush.Command)
}

// NewCommandList builds a list from namespace groups.
// Namespace headers use [yellow]namespace[-] styling and are skipped during navigation.
// onSelect is called when the user selects a command (Enter); nil is safe.
func NewCommandList(groups []drush.NamespaceGroup, onSelect func(cmd *drush.Command)) *CommandList {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Commands")
	list.SetUseStyleTags(true, false)

	headerIndices := make(map[int]bool)
	commandsByIdx := make([]*drush.Command, 0)
	idx := 0

	for _, g := range groups {
		// Add namespace header: styled, non-selectable
		headerText := "[yellow]" + g.Namespace + "[-]"
		list.AddItem(headerText, "", 0, nil)
		headerIndices[idx] = true
		commandsByIdx = append(commandsByIdx, nil)
		idx++

		for i := range g.Commands {
			cmd := &g.Commands[i]
			list.AddItem(cmd.Name, "", 0, nil)
			commandsByIdx = append(commandsByIdx, cmd)
			idx++
		}
	}

	cl := &CommandList{
		List:          list,
		headerIndices: headerIndices,
		commandsByIdx: commandsByIdx,
		onSelect:      onSelect,
	}

	// When selection changes, skip over headers (move to next/prev command)
	list.SetChangedFunc(cl.handleChanged)
	list.SetSelectedFunc(cl.handleSelected)

	return cl
}

// handleChanged runs when the user navigates. If they land on a header, move to nearest command.
func (cl *CommandList) handleChanged(index int, mainText, secondaryText string, shortcut rune) {
	if !cl.headerIndices[index] {
		return
	}
	// Landed on a header — find next command (prefer down, else up)
	next := cl.findNextCommand(index, 1)
	if next < 0 {
		next = cl.findNextCommand(index, -1)
	}
	if next >= 0 {
		cl.SetCurrentItem(next)
	}
}

func (cl *CommandList) findNextCommand(from, direction int) int {
	for i := from + direction; i >= 0 && i < cl.GetItemCount(); i += direction {
		if !cl.headerIndices[i] {
			return i
		}
	}
	return -1
}

// handleSelected runs when the user presses Enter. Ignore headers; call onSelect for commands.
func (cl *CommandList) handleSelected(index int, mainText, secondaryText string, shortcut rune) {
	if cl.headerIndices[index] || cl.onSelect == nil {
		return
	}
	// Find the Command for this index
	cmd := cl.commandAt(index)
	if cmd != nil {
		cl.onSelect(cmd)
	}
}

// commandAt returns the Command at the given list index, or nil if it's a header.
func (cl *CommandList) commandAt(index int) *drush.Command {
	if index < 0 || index >= len(cl.commandsByIdx) {
		return nil
	}
	return cl.commandsByIdx[index]
}
