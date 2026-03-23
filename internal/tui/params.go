package tui

import (
	"fmt"
	"regexp"
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
	error  *tview.TextView
	form   *tview.Form

	// onCancel is called when the user presses Cancel or Esc.
	onCancel func()
	// onRun is called with the command name, positional args, and option flags.
	onRun func(command string, args []string, opts []string)

	// currentHelp tracks which command is displayed so we can read form values.
	currentHelp *drush.CommandHelp
	lastError   string
}

// NewParamsView creates the params pane with placeholder text.
func NewParamsView(onCancel func()) *ParamsView {
	header := tview.NewTextView()
	header.SetDynamicColors(true)

	errorView := tview.NewTextView()
	errorView.SetDynamicColors(true)

	form := tview.NewForm()
	form.SetBorder(false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(errorView, 1, 0, false).
		AddItem(form, 0, 1, true)
	layout.SetBorder(true)
	layout.SetTitle("Parameters")

	pv := &ParamsView{
		layout:   layout,
		header:   header,
		error:    errorView,
		form:     form,
		onCancel: onCancel,
	}

	pv.ShowPlaceholder()
	return pv
}

// ShowPlaceholder resets the pane to its initial state.
func (pv *ParamsView) ShowPlaceholder() {
	pv.clearValidationError()
	pv.currentHelp = nil
	pv.header.SetText("\n  Select a command to see its parameters.")
	pv.form.Clear(true)
}

// ShowError displays an error message in the pane.
func (pv *ParamsView) ShowError(err error) {
	pv.clearValidationError()
	pv.currentHelp = nil
	pv.header.SetText(fmt.Sprintf("\n  [red]Error: %s[-]", err.Error()))
	pv.form.Clear(true)
}

// ShowParams builds a dynamic form from the command's help data.
func (pv *ParamsView) ShowParams(help *drush.CommandHelp) {
	pv.currentHelp = help

	// Header: command name, description, aliases.
	aliasText := ""
	if len(help.Aliases) > 0 {
		aliasText = fmt.Sprintf("\n  Aliases: %s", strings.Join(help.Aliases, ", "))
	}
	pv.header.SetText(fmt.Sprintf("  [yellow]%s[-]  %s%s", help.Name, help.Description, aliasText))
	pv.clearValidationError()

	// Rebuild the form.
	pv.form.Clear(true)

	// Arguments first, sorted by name. Required args get a * prefix.
	argNames := sortedKeys(help.Arguments)
	for _, name := range argNames {
		arg := help.Arguments[name]
		label := name
		required := arg.IsRequired == "1"
		if required {
			label = "* " + label
		}
		fieldName := name
		choices := parseChoices(arg.Description)
		if len(choices) > 0 {
			pv.form.AddDropDown(label, choices, -1, func(text string, index int) {
				if required && text != "" {
					pv.clearValidationErrorFor(fieldName)
				}
			})
		} else {
			pv.form.AddInputField(label, "", 0, nil, func(text string) {
				if required && strings.TrimSpace(text) != "" {
					pv.clearValidationErrorFor(fieldName)
				}
			})
			if desc := arg.Description; desc != "" {
				item := pv.form.GetFormItem(pv.form.GetFormItemCount() - 1)
				if input, ok := item.(*tview.InputField); ok {
					input.SetPlaceholder(desc)
				}
			}
		}
	}

	// Options, sorted by name.
	optNames := sortedKeys(help.Options)
	for _, name := range optNames {
		opt := help.Options[name]
		if opt.AcceptValue == "0" {
			pv.form.AddCheckbox(name, false, nil)
		} else {
			choices := parseChoices(opt.Description)
			if len(choices) > 0 {
				initial := -1
				if len(opt.Defaults) > 0 {
					for i, c := range choices {
						if c == opt.Defaults[0] {
							initial = i
							break
						}
					}
				}
				pv.form.AddDropDown(name, choices, initial, nil)
			} else {
				defaultVal := ""
				if len(opt.Defaults) > 0 {
					defaultVal = opt.Defaults[0]
				}
				pv.form.AddInputField(name, defaultVal, 0, nil, nil)
				if desc := opt.Description; desc != "" {
					item := pv.form.GetFormItem(pv.form.GetFormItemCount() - 1)
					if input, ok := item.(*tview.InputField); ok {
						input.SetPlaceholder(desc)
					}
				}
			}
		}
	}

	// Buttons.
	pv.form.AddButton("Run", func() {
		if pv.onRun == nil || pv.currentHelp == nil {
			return
		}
		if missing := pv.firstMissingRequiredArgument(); missing != "" {
			if pv.lastError != "" {
				// Second press with warning visible — override and run anyway.
				pv.clearValidationError()
			} else {
				pv.showValidationError("Required: " + missing + "  (Run again to skip)")
				return
			}
		} else {
			pv.clearValidationError()
		}
		args, opts := pv.collectValues()
		pv.onRun(pv.currentHelp.Name, args, opts)
	})
	pv.form.AddButton("Cancel", func() {
		pv.ShowPlaceholder()
		if pv.onCancel != nil {
			pv.onCancel()
		}
	})
}

// collectValues reads the current form fields and returns positional args and option flags.
func (pv *ParamsView) collectValues() (args []string, opts []string) {
	help := pv.currentHelp
	if help == nil {
		return
	}

	// Arguments come first in the form (same order as ShowParams).
	argNames := sortedKeys(help.Arguments)
	for _, name := range argNames {
		val := pv.fieldValue(name)
		if val != "" {
			args = append(args, val)
		}
	}

	// Options come after arguments.
	optNames := sortedKeys(help.Options)
	for _, name := range optNames {
		opt := help.Options[name]
		if opt.AcceptValue == "0" {
			item := pv.form.GetFormItemByLabel(name)
			if cb, ok := item.(*tview.Checkbox); ok && cb.IsChecked() {
				opts = append(opts, "--"+name)
			}
		} else {
			val := pv.fieldValue(name)
			if val != "" {
				opts = append(opts, "--"+name+"="+val)
			}
		}
	}

	return
}

func (pv *ParamsView) firstMissingRequiredArgument() string {
	if pv.currentHelp == nil {
		return ""
	}

	for _, name := range sortedKeys(pv.currentHelp.Arguments) {
		arg := pv.currentHelp.Arguments[name]
		if arg.IsRequired != "1" {
			continue
		}
		if strings.TrimSpace(pv.fieldValue(name)) == "" {
			return name
		}
	}

	return ""
}

// fieldValue returns the current text for a form field, checking both the plain
// label and the "* " required prefix. Works for InputField and DropDown items.
func (pv *ParamsView) fieldValue(name string) string {
	item := pv.form.GetFormItemByLabel(name)
	if item == nil {
		item = pv.form.GetFormItemByLabel("* " + name)
	}
	if item == nil {
		return ""
	}
	if input, ok := item.(*tview.InputField); ok {
		return input.GetText()
	}
	if dd, ok := item.(*tview.DropDown); ok {
		_, text := dd.GetCurrentOption()
		return text
	}
	return ""
}

func (pv *ParamsView) showValidationError(message string) {
	pv.lastError = message
	pv.error.SetText("  [red]" + message + "[-]")
}

func (pv *ParamsView) clearValidationError() {
	pv.lastError = ""
	pv.error.SetText("")
}

func (pv *ParamsView) clearValidationErrorFor(fieldName string) {
	if pv.lastError == "Required: "+fieldName {
		pv.clearValidationError()
	}
}

// sortedKeys returns the keys of a map sorted alphabetically.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// choicesPattern matches description text like "Choices: foo, bar." or
// "Available formats: csv,json,list" or "recognized values: a, b, c".
var choicesPattern = regexp.MustCompile(
	`(?i)(?:choices|available\s+\w+|recognized\s+values):\s*([^.]+)`,
)

// parseChoices extracts a list of valid values from a description string.
// Returns nil when no recognizable choices pattern is found.
func parseChoices(desc string) []string {
	m := choicesPattern.FindStringSubmatch(desc)
	if m == nil {
		return nil
	}
	raw := m[1]
	// Normalize " or " and ", " separators to comma.
	raw = strings.ReplaceAll(raw, " or ", ",")
	parts := strings.Split(raw, ",")
	var choices []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			choices = append(choices, p)
		}
	}
	if len(choices) < 2 {
		return nil
	}
	return choices
}
