package blocks

import (
	"fmt"
	"log"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc/event"
	"github.com/libanvl/swager/pkg/ipc/reply"
)

type Tiler struct {
	client  core.Client
	winevts <-chan *event.WindowChange
	opts    *core.Options
}

func init() {
	var _ core.Block = (*Tiler)(nil)
}

func (t *Tiler) Init(client core.Client, sub core.Sub, opts *core.Options) error {
	t.client = client
	t.winevts = sub.Window()
	t.opts = opts
	return nil
}

func (t *Tiler) Configure(args []string) error {
	return nil
}

func (t *Tiler) Run() {
	for evt := range t.winevts {
		if evt.Change != event.FocusWindow {
			if t.opts.Debug {
				t.opts.Log <- "Window change not a focus event"
			}
			continue
		}

		if evt.Container.Type == reply.FloatingConNode {
			if t.opts.Debug {
				t.opts.Log <- "Focused window is floating"
			}
			continue
		}

		t.opts.Log <- fmt.Sprintf("Window: %v, Layout: %v", evt.Container.Name, evt.Container.Layout)

		root, err := t.client.GetTree()
		if err != nil {
			t.opts.Log <- fmt.Sprintf("Fatal error: GetTree failed: %v", err)
			break
		}

		parent := core.FindParent(root, evt.Container.ID)
		if parent == nil {
			if t.opts.Debug {
				t.opts.Log <- "Window has no parent"
			}
			continue
		}

		t.opts.Log <- fmt.Sprintf("Parent: %v, Layout: %v", parent.Type, parent.Layout)

		if parent.Layout == "stacked" || parent.Layout == "tabbed" {
			if t.opts.Debug {
				t.opts.Log <- fmt.Sprintf("Parent layout is excluded: %v", parent.Layout)
			}
			continue
		}

		newlayout := "splith"
		if evt.Container.Rect.Height > evt.Container.Rect.Width {
			newlayout = "splitv"
		}

		if t.opts.Debug {
			t.opts.Log <- fmt.Sprintf("Selecting layout: %v", newlayout)
		}

		if parent.Layout != newlayout {
			s, err := t.client.RunCommand(newlayout)
			if err != nil {
				t.opts.Log <- fmt.Sprintf("Error sending command: %v", err)
			}
			if !(s[0].Success) {
				t.opts.Log <- fmt.Sprintf("sway error: %v", s[0].Error)
			}
		}
	}

	if t.opts.Debug {
		t.opts.Log <- "Tiling channel closed"
	}
}

func (t *Tiler) Close() {
	log.Print("Closing Tiler Block")
}
