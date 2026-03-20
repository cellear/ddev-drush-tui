package tui

import (
	"fmt"

	"github.com/cellear/ddev-drush-tui/internal/drush"
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
	return &OutputView{TextView: tv}
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
