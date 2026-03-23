package tui

import (
	"github.com/cellear/ddev-drush-tui/internal/ddev"
	"github.com/cellear/ddev-drush-tui/internal/drush"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App is the tview application with the three-pane layout.
type App struct {
	app      *tview.Application
	context  *ddev.Context
	commands []drush.NamespaceGroup
	cmdList  *CommandList
	params   *ParamsView
	output   *OutputView
	header   *tview.TextView
	grid     *tview.Grid
}

// NewApp creates the TUI application with header, command list, params form, and output pane.
func NewApp(ctx *ddev.Context, commands []drush.NamespaceGroup) *App {
	app := tview.NewApplication()

	header := tview.NewTextView()
	header.SetDynamicColors(false)
	header.SetText(" ddev-drush-tui | project: " + ctx.ProjectName)

	output := NewOutputView()

	// Build params and command list. Both reference each other, so we
	// create params first with a nil cancel callback, then wire it up.
	params := NewParamsView(nil)

	var cmdList *CommandList
	cmdList = NewCommandList(commands, func(cmd *drush.Command) {
		help, err := drush.Help(cmd.Name)
		if err != nil {
			params.ShowError(err)
			return
		}
		params.ShowParams(help)
		output.SetFocused(false)
		app.SetFocus(params.form)
	})

	params.onCancel = func() {
		output.SetFocused(false)
		app.SetFocus(cmdList.List)
	}

	// Wire the Run button: execute in a goroutine, show output.
	params.onRun = func(command string, args []string, opts []string) {
		go func() {
			app.QueueUpdateDraw(func() {
				output.ShowRunning("ddev drush " + command)
			})
			result, err := drush.Execute(command, args, opts)
			app.QueueUpdateDraw(func() {
				if err != nil {
					output.ShowError(err)
				} else {
					output.ShowResult(result)
				}
			})
		}()
	}

	// Top row: commands (60%) and params (40%)
	panels := tview.NewFlex().
		AddItem(cmdList, 0, 6, true).
		AddItem(params.layout, 0, 4, false)

	// Grid: header, panels, output (output row is focusable for Tab + scrolling)
	grid := tview.NewGrid().
		SetRows(1, -1, 10).
		SetColumns(-1).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(panels, 1, 0, 1, 1, 0, 0, true).
		AddItem(output, 2, 0, 1, 1, 0, 0, true)

	a := &App{
		app:      app,
		context:  ctx,
		commands: commands,
		cmdList:  cmdList,
		params:   params,
		output:   output,
		header:   header,
		grid:     grid,
	}

	app.SetRoot(grid, true)
	app.SetFocus(cmdList.List)

	// isFormFocused returns true when any part of the params form has focus.
	// tview may focus individual form items, not the form itself.
	isFormFocused := func() bool {
		focused := app.GetFocus()
		if focused == params.form {
			return true
		}
		// Check if focus is on a form item (input field, checkbox, button).
		for i := 0; i < params.form.GetFormItemCount(); i++ {
			if focused == params.form.GetFormItem(i) {
				return true
			}
		}
		for i := 0; i < params.form.GetButtonCount(); i++ {
			if focused == params.form.GetButton(i) {
				return true
			}
		}
		return false
	}

	isCommandListFocused := func() bool {
		f := app.GetFocus()
		return f == cmdList.List || f == cmdList
	}

	isOutputFocused := func() bool {
		f := app.GetFocus()
		return f == output || f == output.TextView
	}

	// Arrow keys navigate within the form (down = next field, up = prev field).
	params.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyUp:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		case tcell.KeyTab:
			// Tab cycles to the output pane; use arrows to move between fields.
			output.SetFocused(true)
			app.SetFocus(output)
			return nil
		}
		return event
	})

	// Global key handling.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			if event.Rune() == 'q' && isCommandListFocused() {
				app.Stop()
				return nil
			}
		case tcell.KeyTab:
			if isCommandListFocused() {
				output.SetFocused(false)
				app.SetFocus(params.form)
				return nil
			}
			if isFormFocused() {
				// Handled by form capture when focus is inside the form.
				return event
			}
			if isOutputFocused() {
				output.SetFocused(false)
				app.SetFocus(cmdList.List)
				return nil
			}
			return event
		case tcell.KeyEscape:
			if isFormFocused() {
				params.ShowPlaceholder()
				output.SetFocused(false)
				app.SetFocus(cmdList.List)
				return nil
			}
			if isOutputFocused() {
				output.SetFocused(false)
				app.SetFocus(cmdList.List)
				return nil
			}
			if isCommandListFocused() {
				cmdList.BackToNamespaces()
				return nil
			}
			return event
		}
		return event
	})

	return a
}

// Run starts the TUI and blocks until the user quits.
func (a *App) Run() error {
	return a.app.Run()
}
