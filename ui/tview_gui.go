package ui

import (
        "github.com/rivo/tview"
)

type TviewUI struct {
	app *tview.Application
}

func NewTviewUI() *TviewUI {
	app := tview.NewApplication()
	Configure(app)
	return &TviewUI{
		app: app,
	}
}

func (ui *TviewUI) MainLoop() error {
	return ui.app.Run()
}

func Configure(app *tview.Application) error {
        flex := tview.NewFlex().
                AddItem(tview.NewBox().SetBorder(true).SetTitle("Left (1/2 x width of Top)"), 0, 1, false).
                AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
                AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	app.SetRoot(flex, true)
	return nil
}
