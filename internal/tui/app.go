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
		app.SetFocus(params.form)
	})

	params.onCancel = func() {
		app.SetFocus(cmdList.List)
	}

	// Top row: commands (60%) and params (40%)
	panels := tview.NewFlex().
		AddItem(cmdList, 0, 6, true).
		AddItem(params.layout, 0, 4, false)

	// Grid: header, panels, output
	grid := tview.NewGrid().
		SetRows(1, -1, 5).
		SetColumns(-1).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(panels, 1, 0, 1, 1, 0, 0, true).
		AddItem(output, 2, 0, 1, 1, 0, 0, false)

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
	app.SetFocus(cmdList)

	// Global key handling: q quits from command list, Tab cycles focus, Esc returns to command list.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			if event.Rune() == 'q' && app.GetFocus() == cmdList.List {
				app.Stop()
				return nil
			}
		case tcell.KeyTab:
			// Toggle focus between command list and params form.
			focused := app.GetFocus()
			if focused == cmdList.List {
				app.SetFocus(params.form)
			} else {
				app.SetFocus(cmdList.List)
			}
			return nil
		case tcell.KeyEscape:
			// Esc always returns to command list.
			app.SetFocus(cmdList.List)
			return nil
		}
		return event
	})

	return a
}

// Run starts the TUI and blocks until the user quits.
func (a *App) Run() error {
	return a.app.Run()
}
