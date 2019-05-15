package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"jcheng/grs/script"
	"sync"
	"time"
)

func colorI(s script.Indexstat) string {
	if s == script.INDEX_UNMODIFIED {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}

func colorIGrs(s script.Indexstat) string {
	if s == script.INDEX_UNMODIFIED {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\033[31m%v\033[0m", s)
}

func colorB(s script.Branchstat) string {
	if s == script.BRANCH_UPTODATE {
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
	return nil
}

// === START: CliUI implementation ===
type CliUI interface {
	DoneSender() <-chan struct{}
	MainLoop() error
	DrawGrs(repo []script.GrsRepo)
	Close()
}

type ConsoleUI struct {
	gui      *gocui.Gui
	done     chan struct{}
	doneLock sync.Mutex
}

var _consoleUIImpl CliUI = &ConsoleUI{}

func NewConsoleUI() (*ConsoleUI, error) {
	gui, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}

	gui.SetManagerFunc(_layout)

	consoleUI := &ConsoleUI{
		gui:  gui,
		done: make(chan struct{}),
	}

	if err := consoleUI.initKeyBindings(); err != nil {
		return nil, err
	}

	return consoleUI, nil
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
	return nil
}

func (consoleUI *ConsoleUI) DoneSender() <-chan struct{} {
	return consoleUI.done
}

func (consoleUI *ConsoleUI) MainLoop() error {
	return consoleUI.gui.MainLoop()
}

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

		for _, repo := range repos {
			pushIndicator := ""
			if repo.IsPushAllowed() {
				pushIndicator = "\033[32m⯅\033[0m"
			}

			line := fmt.Sprintf("repo [%v]%v status is %v, %v, %v.",
				repo.GetLocal(), pushIndicator, colorBGrs(repo.GetStats().Branch), colorIGrs(repo.GetStats().Index), repo.GetStats().CommitTime)
			if _, err := fmt.Fprintln(v, line); err != nil {
				return err
			}
		}
		return nil
	})
}

func (consoleUI *ConsoleUI) Close() {
	consoleUI.gui.Close()
}

type PrintUI struct {
	done chan struct{}
}

func NewPrintUI() (*PrintUI, error) {
	return &PrintUI{
		done: make(chan struct{}),
	}, nil
}

func (printUI *PrintUI) DoneSender() <-chan struct{} {
	return printUI.done
}

func (printUI *PrintUI) MainLoop() error {
	<-printUI.done
	return nil
}

func (printUI *PrintUI) DrawGrs(repos []script.GrsRepo) {
	fmt.Print("\033[2J\033[H")
	fmt.Println(time.Now().Format("=== Jan _2 3:04PM MST ==="))

	for _, repo := range repos {
		fmt.Printf("repo [%v] status IS %v, %v, %v.\n",
			repo.GetLocal(), repo.GetStats().Branch, repo.GetStats().Index, repo.GetStats().CommitTime)
	}
}

func (printUI *PrintUI) Close() {
	close(printUI.done)
}

var _printUIImpl CliUI = &PrintUI{}

// === END: CliUI implementation ===
