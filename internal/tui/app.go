package tui

import (
	"github.com/cellear/ddev-drush-tui/internal/ddev"
	"github.com/cellear/ddev-drush-tui/internal/drush"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App is the tview application with the three-pane layout.
type App struct {
	app       *tview.Application
	context   *ddev.Context
	commands  []drush.NamespaceGroup
	cmdList   *CommandList
	params    *ParamsView
	output    *OutputView
	header    *tview.TextView
	grid      *tview.Grid
}

// NewApp creates the TUI application with header, command list, params placeholder, and output placeholder.
func NewApp(ctx *ddev.Context, commands []drush.NamespaceGroup) *App {
	header := tview.NewTextView()
	header.SetDynamicColors(false)
	header.SetText(" ddev-drush-tui | project: " + ctx.ProjectName)

	// Command list with no-op callback for now (S2 will wire to params)
	cmdList := NewCommandList(commands, func(cmd *drush.Command) {
		// S2: load params for cmd
		_ = cmd
	})

	params := NewParamsView()
	output := NewOutputView()

	// Top row: commands (60%) and params (40%)
	panels := tview.NewFlex().
		AddItem(cmdList, 0, 6, true).  // 60% = 6/10
		AddItem(params, 0, 4, false)   // 40% = 4/10

	// Grid: header, panels, output
	grid := tview.NewGrid().
		SetRows(1, -1, 5).
		SetColumns(-1).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(panels, 1, 0, 1, 1, 0, 0, true).
		AddItem(output, 2, 0, 1, 1, 0, 0, false)

	app := tview.NewApplication()

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

	// q quits when command list is focused
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			if app.GetFocus() == cmdList {
				app.Stop()
				return nil
			}
		}
		return event
	})

	return a
}

// Run starts the TUI and blocks until the user quits.
func (a *App) Run() error {
	return a.app.Run()
}
