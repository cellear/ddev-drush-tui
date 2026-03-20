package tui

import (
	"github.com/rivo/tview"
)

// OutputView is the bottom pane placeholder showing command output.
// S3 will replace this with scrollable output from executed commands.
type OutputView struct {
	*tview.TextView
}

// NewOutputView creates the output placeholder.
func NewOutputView() *OutputView {
	tv := tview.NewTextView()
	tv.SetBorder(true)
	tv.SetTitle("Output")
	tv.SetText("Output will appear here after running a command.")
	tv.SetDynamicColors(false)
	return &OutputView{TextView: tv}
}
