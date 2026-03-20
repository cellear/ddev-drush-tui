package tui

import (
	"github.com/rivo/tview"
)

// ParamsView is the right pane placeholder showing parameter form.
// S2 will replace this with a dynamic form for the selected command.
type ParamsView struct {
	*tview.TextView
}

// NewParamsView creates the params placeholder.
func NewParamsView() *ParamsView {
	tv := tview.NewTextView()
	tv.SetBorder(true)
	tv.SetTitle("Parameters")
	tv.SetText("Select a command to see its parameters.")
	tv.SetDynamicColors(false)
	return &ParamsView{TextView: tv}
}
