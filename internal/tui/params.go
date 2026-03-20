package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cellear/ddev-drush-tui/internal/drush"
	"github.com/rivo/tview"
)

// ParamsView is the right pane showing a dynamic parameter form
// for the selected Drush command.
type ParamsView struct {
	// layout holds the header text and the form in a vertical flex.
	layout *tview.Flex
	header *tview.TextView
	form   *tview.Form

	// onCancel is called when the user presses Cancel or Esc.
	onCancel func()
}

// NewParamsView creates the params pane with placeholder text.
func NewParamsView(onCancel func()) *ParamsView {
	header := tview.NewTextView()
	header.SetDynamicColors(true)

	form := tview.NewForm()
	form.SetBorder(false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(form, 0, 1, true)
	layout.SetBorder(true)
	layout.SetTitle("Parameters")

	pv := &ParamsView{
		layout:   layout,
		header:   header,
		form:     form,
		onCancel: onCancel,
	}

	pv.ShowPlaceholder()
	return pv
}

// ShowPlaceholder resets the pane to its initial state.
func (pv *ParamsView) ShowPlaceholder() {
	pv.header.SetText("\n  Select a command to see its parameters.")
	pv.form.Clear(true)
}

// ShowError displays an error message in the pane.
func (pv *ParamsView) ShowError(err error) {
	pv.header.SetText(fmt.Sprintf("\n  [red]Error: %s[-]", err.Error()))
	pv.form.Clear(true)
}

// ShowParams builds a dynamic form from the command's help data.
func (pv *ParamsView) ShowParams(help *drush.CommandHelp) {
	// Header: command name, description, aliases.
	aliasText := ""
	if len(help.Aliases) > 0 {
		aliasText = fmt.Sprintf("\n  Aliases: %s", strings.Join(help.Aliases, ", "))
	}
	pv.header.SetText(fmt.Sprintf("  [yellow]%s[-]  %s%s", help.Name, help.Description, aliasText))

	// Rebuild the form.
	pv.form.Clear(true)

	// Arguments first, sorted by name. Required args get a * prefix.
	argNames := sortedKeys(help.Arguments)
	for _, name := range argNames {
		arg := help.Arguments[name]
		label := name
		if arg.IsRequired == "1" {
			label = "* " + label
		}
		pv.form.AddInputField(label, "", 0, nil, nil)
	}

	// Options, sorted by name.
	optNames := sortedKeys(help.Options)
	for _, name := range optNames {
		opt := help.Options[name]
		if opt.AcceptValue == "0" {
			// Boolean flag — render as checkbox.
			pv.form.AddCheckbox(name, false, nil)
		} else {
			// Value option — render as input field with default if available.
			defaultVal := ""
			if len(opt.Defaults) > 0 {
				defaultVal = opt.Defaults[0]
			}
			pv.form.AddInputField(name, defaultVal, 0, nil, nil)
		}
	}

	// Buttons.
	pv.form.AddButton("Run", func() {
		// S3 will wire this to drush.Execute().
	})
	pv.form.AddButton("Cancel", func() {
		pv.ShowPlaceholder()
		if pv.onCancel != nil {
			pv.onCancel()
		}
	})
}

// sortedKeys returns the keys of a map sorted alphabetically.
// Works for both Argument and Option maps via a generic-free approach.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
