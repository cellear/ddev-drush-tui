package tui

import (
	"fmt"

	"github.com/cellear/ddev-drush-tui/internal/drush"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// OutputView is the bottom pane showing command output.
type OutputView struct {
	*tview.TextView
}

// NewOutputView creates the output pane with placeholder text.
func NewOutputView() *OutputView {
	tv := tview.NewTextView()
	tv.SetBorder(true)
	tv.SetTitle("Output")
	tv.SetText("Output will appear here after running a command.")
	tv.SetDynamicColors(true)
	tv.SetScrollable(true)

	ov := &OutputView{TextView: tv}

	// less-style extras: Space / Page Down = page down; b / Page Up = page up.
	// Arrow keys, Home/End, g/G are handled by tview.TextView.
	tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case ' ':
				return tcell.NewEventKey(tcell.KeyPgDn, 0, tcell.ModNone)
			case 'b':
				return tcell.NewEventKey(tcell.KeyPgUp, 0, tcell.ModNone)
			}
		}
		return event
	})

	return ov
}

// SetFocused updates the border color so the user can see when this pane has focus.
func (ov *OutputView) SetFocused(focused bool) {
	if focused {
		ov.SetBorderColor(tcell.ColorYellow)
	} else {
		ov.SetBorderColor(tview.Styles.BorderColor)
	}
}

// ShowRunning displays a "running..." message while a command executes.
func (ov *OutputView) ShowRunning(command string) {
	ov.SetText(fmt.Sprintf("$ %s\n\n[yellow]Running...[-]", command))
}

// ShowResult displays the command output. Errors are highlighted.
func (ov *OutputView) ShowResult(result *drush.Result) {
	header := fmt.Sprintf("$ %s\n\n", result.Command)
	if result.ExitCode != 0 {
		ov.SetText(fmt.Sprintf("%s%s\n\n[red]Exit code: %d[-]", header, result.Output, result.ExitCode))
	} else {
		ov.SetText(header + result.Output)
	}
	ov.ScrollToBeginning()
}

// ShowError displays an execution error.
func (ov *OutputView) ShowError(err error) {
	ov.SetText(fmt.Sprintf("[red]Error: %s[-]", err.Error()))
}
