package tui

import (
	"fmt"
	"strings"

	"github.com/cellear/ddev-drush-tui/internal/drush"
	"github.com/gdamore/tcell/v2"
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
	*tview.Flex
	list             *tview.List
	filter           *tview.InputField
	groups           []drush.NamespaceGroup
	level            int
	currentNamespace string
	filterActive     bool
	filterText       string
	currentItems     []commandListItem
	onSelect         func(cmd *drush.Command)
}

type commandListItem struct {
	label     string
	namespace string
	command   *drush.Command
	isBack    bool
}

// NewCommandList builds a namespace-first list. onSelect is called when the
// user chooses a command at level 2; nil is safe.
func NewCommandList(groups []drush.NamespaceGroup, onSelect func(cmd *drush.Command)) *CommandList {
	list := tview.NewList()
	list.SetBorder(false)
	list.SetUseStyleTags(true, false)

	filter := tview.NewInputField()
	filter.SetLabel("/ ")
	filter.SetFieldWidth(0)

	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true)
	container.AddItem(filter, 0, 0, false)
	container.AddItem(list, 0, 1, true)

	cl := &CommandList{
		Flex:     container,
		list:     list,
		filter:   filter,
		groups:   groups,
		level:    cmdLevelNamespaces,
		onSelect: onSelect,
	}
	filter.SetChangedFunc(func(text string) {
		cl.filterText = text
		cl.refresh()
	})
	filter.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			cl.HideFilter()
		case tcell.KeyEnter:
			cl.selectCurrent()
			cl.HideFilter()
		}
	})
	cl.showNamespaces()
	list.SetSelectedFunc(cl.handleSelected)
	return cl
}

func (cl *CommandList) showNamespaces() {
	cl.level = cmdLevelNamespaces
	cl.currentNamespace = ""
	cl.refresh()
}

func (cl *CommandList) showCommands(namespace string) {
	cl.level = cmdLevelCommands
	cl.currentNamespace = namespace
	cl.refresh()
}

func (cl *CommandList) handleSelected(index int, mainText, secondaryText string, shortcut rune) {
	if index < 0 || index >= len(cl.currentItems) {
		return
	}

	item := cl.currentItems[index]
	if cl.level == cmdLevelNamespaces {
		cl.showCommands(item.namespace)
		return
	}
	if item.isBack {
		cl.showNamespaces()
		return
	}
	if item.command != nil && cl.onSelect != nil {
		cl.onSelect(item.command)
	}
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

// ShowFilter reveals the inline filter input for the current level.
func (cl *CommandList) ShowFilter() {
	cl.filterActive = true
	cl.filter.SetText(cl.filterText)
	cl.ResizeItem(cl.filter, 1, 0)
}

// HideFilter clears the filter and restores the full list for the current level.
func (cl *CommandList) HideFilter() {
	cl.filterActive = false
	cl.filterText = ""
	cl.filter.SetText("")
	cl.ResizeItem(cl.filter, 0, 0)
	cl.refresh()
}

// IsFilterFocused reports whether the command pane's filter input currently has focus.
func (cl *CommandList) IsFilterFocused(focused tview.Primitive) bool {
	return focused == cl.filter
}

// ListPrimitive returns the underlying list for focus management.
func (cl *CommandList) ListPrimitive() *tview.List {
	return cl.list
}

// FilterPrimitive returns the inline filter input.
func (cl *CommandList) FilterPrimitive() *tview.InputField {
	return cl.filter
}

func (cl *CommandList) selectCurrent() {
	index := cl.list.GetCurrentItem()
	main, secondary := cl.list.GetItemText(index)
	cl.handleSelected(index, main, secondary, 0)
}

func (cl *CommandList) refresh() {
	cl.list.Clear()
	cl.currentItems = cl.filteredItems()
	cl.setTitle()
	for _, item := range cl.currentItems {
		cl.list.AddItem(item.label, "", 0, nil)
	}
	if cl.list.GetItemCount() > 0 {
		cl.list.SetCurrentItem(0)
	}
}

func (cl *CommandList) setTitle() {
	if cl.level == cmdLevelCommands {
		cl.SetTitle("Commands > " + cl.currentNamespace)
		return
	}
	cl.SetTitle("Commands")
}

func (cl *CommandList) filteredItems() []commandListItem {
	query := strings.TrimSpace(strings.ToLower(cl.filterText))
	if cl.level == cmdLevelCommands {
		return cl.filteredCommandItems(query)
	}
	return cl.filteredNamespaceItems(query)
}

func (cl *CommandList) filteredNamespaceItems(query string) []commandListItem {
	items := make([]commandListItem, 0, len(cl.groups))
	for i := range cl.groups {
		group := cl.groups[i]
		if query != "" && !strings.Contains(strings.ToLower(group.Namespace), query) {
			continue
		}
		items = append(items, commandListItem{
			label:     fmt.Sprintf("%s  (%d)", group.Namespace, len(group.Commands)),
			namespace: group.Namespace,
		})
	}
	return items
}

func (cl *CommandList) filteredCommandItems(query string) []commandListItem {
	items := []commandListItem{{
		label:  "[yellow]← Back[-]",
		isBack: true,
	}}

	for i := range cl.groups {
		if cl.groups[i].Namespace != cl.currentNamespace {
			continue
		}
		for j := range cl.groups[i].Commands {
			cmd := &cl.groups[i].Commands[j]
			if query != "" && !strings.Contains(strings.ToLower(cmd.Name), query) {
				continue
			}
			items = append(items, commandListItem{
				label:   cmd.Name,
				command: cmd,
			})
		}
		break
	}
	return items
}
