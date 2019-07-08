package ui

import (
        "github.com/rivo/tview"
	"github.com/gdamore/tcell"
)

// TviewUI renders the UI using the tview framework. The CliUI interface requires
// 	DoneSender() <-chan struct{}
//	EventSender() <-chan UiEvent
//	MainLoop() error
//	DrawGrs(repo []script.GrsRepo)
//	Close()
//
//  Obvious candidates are:
//      MainLoop() -> Application.Run()
//      Stop() -> Application.Stop()
//
//  The DrawGrs function should clear the Application object and reconfiguire it from scratch.
//  This is a little inefficient, as we get rid of reusable objects like the Flex and Table
//  object. However, since we only refresh once every few minutes, it is not worth trading
//  off simplicity for efficiency.
//  
type TviewUI struct {
	app *tview.Application
	flex *tview.Flex
	bottomPane tview.Primitive
}

func NewTviewUI() *TviewUI {
	app := tview.NewApplication()
	flex, bottomPane := Configure(app)
	return &TviewUI{
		app: app,
		flex: flex,
		bottomPane: bottomPane,
	}
}

func (ui *TviewUI) MainLoop() error {
	return ui.app.Run()
}

func Configure(app *tview.Application) (*tview.Flex, tview.Primitive) {
	bottomPane := NewTviewExample()
        flex := tview.NewFlex().SetDirection(tview.FlexRow).
                AddItem(tview.NewBox().SetBorder(true).SetTitle("A"), 0, 1, false).
		AddItem(bottomPane, 5, 1, false)
	app.SetRoot(flex, true)
	return flex, bottomPane
}

type TviewExample struct {
	textView *tview.TextView
}

func NewTviewExample() *TviewExample {
	retval := TviewExample{
		textView: tview.NewTextView().SetText("B (5 rows)"),
	}
	return &retval
}

func (s *TviewExample) Draw(screen tcell.Screen) {
	s.textView.Draw(screen)
}

func (s *TviewExample) GetRect() (int, int, int, int) {
	return s.textView.GetRect()
}

func (s *TviewExample) SetRect(x, y, width, height int) {
	s.textView.SetRect(x, y, width, height)
}

func (s *TviewExample) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return s.textView.InputHandler()
}

func (s *TviewExample) Focus(delegate func(p tview.Primitive)) {
	s.textView.Focus(delegate)
}

func (s *TviewExample) Blur() {
	s.textView.Blur()
}

func (s *TviewExample) GetFocusable() tview.Focusable {
	return s.textView.GetFocusable()
}
