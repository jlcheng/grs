package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"jcheng/grs/script"
	"log"
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
}

func NewCuiGUI() *CuiGUI {
	return &CuiGUI{}
}

func (c *CuiGUI) Init() error {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return err
	}

	c.cui = g
	g.SetManagerFunc(_cuiManager)
	return c.initKeyBindings()
}

func (c *CuiGUI) initKeyBindings() error {
	quitFunc := func(_ *gocui.Gui, _ *gocui.View) error {
		c.stopped = true
		return gocui.ErrQuit
	}
	if err := c.cui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, quitFunc); err != nil {
		return err
	}
	return nil
}

func (c *CuiGUI) Run(repos []script.Repo) {
	c.cui.Update(func(g *gocui.Gui) error {
		v, err := g.View("main")
		if err != nil {
			return err
		}
		v.Clear()

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
	c.cui.Close()
}

func _cuiManager(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Grs"
	}
	return nil
}
// === END: CUI implementation ===