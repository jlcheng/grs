package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"jcheng/grs/script"
)

// TviewUI renders the UI using the tview framework.
//
// The DrawGrs clears the Application object and reconfiguire it from
// scratch.  This is inefficient, as we get rid of reusable objects
// such as Flex and Table. However, since we only refresh once every
// few minutes, it is a worthwhile trade-off of efficiency for
// readbility.
type TviewUI struct {
	app     *tview.Application
	eventCh chan UiEvent
}

func NewTviewUI() *TviewUI {
	app := tview.NewApplication()

	return &TviewUI{
		app:     app,
		eventCh: make(chan UiEvent),
	}
}

func (ui *TviewUI) MainLoop() error {
	return ui.app.Run()
}

func (ui *TviewUI) Close() {
	ui.app.Stop()
}

func (ui *TviewUI) DrawGrs(repo []script.GrsRepo) {
	bottomPane := NewTviewExample()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Grs"), 0, 1, false).
		AddItem(bottomPane, 5, 1, false)
	ui.app.SetRoot(flex, true)
	ui.app.Draw()
}

func (ui *TviewUI) EventSender() <-chan UiEvent {
	return ui.eventCh
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
