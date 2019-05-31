package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"jcheng/grs/script"
	"strings"
	"sync"
	"time"
)

func colorIGrs(s script.Indexstat) string {
	if s == script.INDEX_UNMODIFIED {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}

func colorBGrs(s script.Branchstat) string {
	if s == script.BRANCH_UPTODATE {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}

func _layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Grs"
		if _, err := fmt.Fprintln(v, "Fetching repo data..."); err != nil {
			return err
		}
	}

	// TODO: Determine this number from the number of repos in the dashboard
	minDashboard := 7
	errorTop := maxY - 12
	errorBottom := maxY - 2
	errorLeft := 1
	errorRight := maxX - 2
	if maxY > minDashboard {

		if errorTop < minDashboard-2 {
			errorTop = minDashboard - 2
		}
		if errorBottom < minDashboard-1 {
			errorBottom = minDashboard - 1
		}

		if v, err := g.SetView("errors", errorLeft, errorTop, errorRight, errorBottom); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Title = "Errors"
			v.Frame = true
			v.Editable = false
			v.Highlight = false
			v.SelFgColor = gocui.ColorCyan
			v.Overwrite = false
			if _, err := g.SetCurrentView("errors"); err != nil {
				return err
			}
		}
	}

	return nil
}

type CliUI interface {
	DoneSender() <-chan struct{}
	EventSender() <-chan UiEvent
	MainLoop() error
	DrawGrs(repo []script.GrsRepo)
	Close()
}

// ConsoleUI is a prettier and more powerful UI implementation
type ConsoleUI struct {
	gui      *gocui.Gui
	done     chan struct{}
	doneLock sync.Mutex
	eventCh  chan UiEvent
}

// NewConsoleUI creates a ConsoleUI and initialize its UI layout and keybindings.
func NewConsoleUI() (*ConsoleUI, error) {
	gui, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}

	gui.SetManagerFunc(_layout)
	gui.Cursor = true

	consoleUI := &ConsoleUI{
		gui:     gui,
		done:    make(chan struct{}),
		eventCh: make(chan UiEvent),
	}

	if err := consoleUI.initKeyBindings(); err != nil {
		return nil, err
	}

	return consoleUI, nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (consoleUI *ConsoleUI) initKeyBindings() error {
	quitFunc := func(_ *gocui.Gui, _ *gocui.View) error {
		consoleUI.doneLock.Lock()
		defer consoleUI.doneLock.Unlock()
		if consoleUI.done == nil {
			return nil
		}
		close(consoleUI.done)
		consoleUI.done = nil
		return gocui.ErrQuit
	}
	if err := consoleUI.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitFunc); err != nil {
		return err
	}

	if err := consoleUI.gui.SetKeybinding("errors", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := consoleUI.gui.SetKeybinding("errors", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	refreshFunc := func(g *gocui.Gui, _ *gocui.View) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}

		v.Title = "Refreshing"
		// Careful: If the event queue is full, the refresh event will be lost.
		select {
		case consoleUI.eventCh <- EVENT_REFRESH:
		default:
		}

		return nil
	}
	if err := consoleUI.gui.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, refreshFunc); err != nil {
		return err
	}
	return nil
}

// DoneSender returns a channel that blocks until the ConsoleUI is closed
func (consoleUI *ConsoleUI) DoneSender() <-chan struct{} {
	return consoleUI.done
}

// MainLoop blocks until the UI loop is complete
func (consoleUI *ConsoleUI) MainLoop() error {
	return consoleUI.gui.MainLoop()
}

// DrawGrs enqueues a draw operation in the UI's rendering pipeline
func (consoleUI *ConsoleUI) DrawGrs(repos []script.GrsRepo) {
	if consoleUI.done == nil {
		return
	}

	consoleUI.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		v.Clear()
		var timestr = time.Now().Format("[Jan _2 3:04:05PM PST]")
		v.Title = fmt.Sprintf("Grs %s", timestr)

		errView, err := g.View("errors")
		if err != nil {
			return err
		}
		errView.Clear()

		for _, repo := range repos {
			pushIndicator := ""
			if repo.IsPushAllowed() {
				pushIndicator = "\033[32m⯅\033[0m"
			}

			errorIndicator := " "
			if repo.GetError() != nil {
				errorIndicator = "\033[31;47m!\033[0m"
				errorDetails := strings.Trim(repo.GetError().Error(), "\n")
				errorMessage := fmt.Sprintf("\033[32m%v\033[0m\n%v\n", repo.GetLocal(), errorDetails)
				fmt.Fprintln(errView, errorMessage)
			}

			line := fmt.Sprintf("%vrepo [%v]%v status is %v, %v, %v.",
				errorIndicator, repo.GetLocal(), pushIndicator, colorBGrs(repo.GetStats().Branch),
				colorIGrs(repo.GetStats().Index), repo.GetStats().CommitTime)

			// Writes any error messages to the error view
			if _, err := fmt.Fprintln(v, line); err != nil {
				return err
			}
		}

		return nil
	})
}

// Close releases resources used by this ConsoleUI instance
func (consoleUI *ConsoleUI) Close() {
	consoleUI.gui.Close()
}

// EventSender returns a channel that someone can poll for UI events
func (consoleUI *ConsoleUI) EventSender() <-chan UiEvent {
	return consoleUI.eventCh
}

// PrintUI is the simpler and less useful implementation of CliUI
type PrintUI struct {
	done    chan struct{}
	eventCh <-chan UiEvent
}

// NewPrintUI returns a PrintUI instance
func NewPrintUI() (*PrintUI, error) {
	return &PrintUI{
		done:    make(chan struct{}),
		eventCh: make(chan UiEvent),
	}, nil
}

// DoneSender returns a channel that blocks until the PrintUI is closed
func (printUI *PrintUI) DoneSender() <-chan struct{} {
	return printUI.done
}

// MainLoop blocks until the UI loop is complete
func (printUI *PrintUI) MainLoop() error {
	<-printUI.done
	return nil
}

// DrawGrs blocks while it draws the state of the given GrsRepo array to the screen
func (printUI *PrintUI) DrawGrs(repos []script.GrsRepo) {
	fmt.Print("\033[2J\033[H")
	fmt.Println(time.Now().Format("=== Jan _2 3:04PM MST ==="))

	for _, repo := range repos {
		fmt.Printf("repo [%v] status IS %v, %v, %v.\n",
			repo.GetLocal(), repo.GetStats().Branch, repo.GetStats().Index, repo.GetStats().CommitTime)
	}
}

// Close releases resources used by this PrintUI instance
func (printUI *PrintUI) Close() {
	close(printUI.done)
}

// EventSender returns a channel that always blocks, as the PrintUI object is too simple to generate UI events
func (printUI *PrintUI) EventSender() <-chan UiEvent {
	return printUI.eventCh
}
