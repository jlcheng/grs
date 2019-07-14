package ui

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"jcheng/grs/base"
	"jcheng/grs"
	"strings"
	"time"
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
	box     *tview.Box
	eventCh chan UiEvent
}

func NewTviewUI() *TviewUI {
	app := tview.NewApplication()
	ui := &TviewUI{
		app:     app,
		box:     nil,
		eventCh: make(chan UiEvent),
	}
	
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
		}
		if event.Key() == tcell.KeyCtrlR {
			if ui.box != nil {
				ui.box.SetTitle("Refreshing")
			}
			// Careful: If the event queue is full, the refresh event will be lost.
			select {
			case ui.eventCh <- EVENT_REFRESH:
			default:
			}
		}
		return event
	})

	ui.DrawGrs(make([]grs.GrsRepo,0))
	ui.box.SetTitle("Starting Grs...")
	
	return ui
}

func (ui *TviewUI) MainLoop() error {
	return ui.app.Run()
}

func (ui *TviewUI) Close() {
	ui.app.Stop()
}

func (ui *TviewUI) DrawGrs(repos []grs.GrsRepo) {
	errorCount, textView := errorView(repos)
	reposPane := repositoryTable(repos)
	ui.box = reposPane.Box
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(reposPane, 0, 1, false)

	if errorCount != 0 {
		flex.AddItem(textView, errorCount+2, 1, true)
		ui.app.SetFocus(textView)
	}
	ui.app.SetRoot(flex, true)
	ui.app.Draw()
}

func (ui *TviewUI) EventSender() <-chan UiEvent {
	return ui.eventCh
}

func repositoryTable(repos []grs.GrsRepo) *tview.Table {
	timeStr := time.Now().Format("Jan _2 3:04:05PM PST")
	title := fmt.Sprintf("Grs (%v) [%v]", base.Version(), timeStr)
	
	table := tview.NewTable()
	table.SetCell(0, 0, tview.NewTableCell("Push").SetTextColor(tcell.ColorGreen))
	table.SetCell(0, 1, tview.NewTableCell("Path").SetTextColor(tcell.ColorGreen))
	table.SetCell(0, 2, tview.NewTableCell("Local").SetTextColor(tcell.ColorGreen))
	table.SetCell(0, 3, tview.NewTableCell("Branch").SetTextColor(tcell.ColorGreen))
	table.SetCell(0, 4, tview.NewTableCell("Index").SetTextColor(tcell.ColorGreen))
	
	for i := 0; i < len(repos); i++ {
		repo := repos[i]
		pushIndicator := ""
		if repo.IsPushAllowed() {
			pushIndicator = "↑↓"
		}
		table.SetCell(i+1, 0,
			tview.NewTableCell(pushIndicator).SetTextColor(tcell.ColorGreen))
		
		table.SetCell(i+1, 1,
			tview.NewTableCell(repos[i].GetLocal()))

		color := tcell.ColorWhite
		text := repos[i].GetStats().Dir.String()
		if repos[i].GetStats().Dir != grs.GRSDIR_VALID {
			color = tcell.ColorRed
		}
		table.SetCell(i+1, 2, tview.NewTableCell(text).SetTextColor(color))

		color = tcell.ColorWhite
		text = repos[i].GetStats().Branch.String()
		if repos[i].GetStats().Branch != grs.BRANCH_UPTODATE {
			color = tcell.ColorRed
		}
		table.SetCell(i+1, 3, tview.NewTableCell(text).SetTextColor(color))

		color = tcell.ColorWhite
		text = repos[i].GetStats().Index.String()
		if repos[i].GetStats().Index != grs.INDEX_UNMODIFIED {
			color = tcell.ColorRed
		}
		table.SetCell(i+1, 4, tview.NewTableCell(text).SetTextColor(color))
	}

	table.Box.SetBorder(true)
	table.Box.SetTitle(title)
	return table
}

func errorView(repos []grs.GrsRepo) (int, *tview.TextView) {
	errorCount := 0
	errorMsg := ""
	for _, repo := range repos {
		if repo.GetError() != nil {
			errorCount++
			errorMsg = errorMsg + repo.GetError().Error()
		}
	}
	
	textView := tview.NewTextView()
	textView.SetBorder(true).SetTitle("errors")
	textView.ScrollToEnd()
	errorMsg = strings.Trim(errorMsg, "\n")
	fmt.Fprint(textView, errorMsg)
	return errorCount, textView
}
