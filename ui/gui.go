package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"jcheng/grs/script"
	"log"
	"sync"
	"time"
)

type AnsiGUI struct {
	clr bool // if true, clears screen before each iteration
}

func NewGUI(clr bool) AnsiGUI {
	return AnsiGUI{
		clr: clr,
	}
}

func (gui *AnsiGUI) Run(repos []script.Repo) {
	// setup/clear screen
	if gui.clr {
		fmt.Print("\033[2J\033[H")
	}
	fmt.Println(time.Now().Format("=== Jan _2 3:04PM MST ==="))

	for _, repo := range repos {
		fmt.Printf("repo [%v] status IS %v, %v, %v.\n",
			repo.Path, colorB(repo.Branch), colorI(repo.Index), repo.CommitTime)
	}
}

func colorI(s script.Indexstat) string {
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


// === START: CUI implementation ===
type CuiGUI struct {
	cui *gocui.Gui
	stopped bool
	done chan struct{}
	doneLock sync.Mutex
}

func NewCuiGUI() *CuiGUI {
	return &CuiGUI{}
}

func (c *CuiGUI) Init() error {
	g, err := gocui.NewGui(gocui.Output256)
	g.ASCII = false

	if err != nil {
		return err
	}

	c.cui = g
	g.SetManagerFunc(_layout)
	if err := c.initKeyBindings(); err != nil {
		return err
	}

	c.doneLock.Lock()
	c.done = make(chan struct{})
	c.doneLock.Unlock()

	return nil
}

func (c *CuiGUI) initKeyBindings() error {
	quitFunc := func(_ *gocui.Gui, _ *gocui.View) error {
		c.stopped = true

		c.doneLock.Lock()
		close(c.done)
		c.done = nil
		c.doneLock.Unlock()

		return gocui.ErrQuit
	}
	if err := c.cui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitFunc); err != nil {
		return err
	}
	return nil
}

func (c *CuiGUI) Run(repos []script.Repo) {
	if c.stopped {
		return
	}

	c.cui.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		v.Clear()
		var time = time.Now().Format("[Jan _2 3:04PM PST]")
		v.Title = fmt.Sprintf("Grs %s", time)

		for _, repo := range repos {
			line := fmt.Sprintf("repo [%v] status IS %v, %v, %v.",
				repo.Path, colorB(repo.Branch), colorI(repo.Index), repo.CommitTime)
			fmt.Fprintln(v, line)
		}
		return nil
	})
}

func (c *CuiGUI) MainLoop() error {
	if err := c.cui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
	return nil
}

func (c *CuiGUI) Close() {
	c.doneLock.Lock()
	if c.done != nil {
		close(c.done)
	}
	c.doneLock.Unlock()

	c.cui.Close()
}

func _layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Grs"
		fmt.Fprintln(v, "Fetching repo data...")
	}
	return nil
}
// === END: CUI implementation ===