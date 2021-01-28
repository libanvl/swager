package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type Tiler struct {
	client  core.Client
	winevts <-chan *ipc.WindowChange
	opts    *core.Options
}

func init() {
	var _ core.Block = (*Tiler)(nil)
}

func (t *Tiler) Init(client core.Client, sub core.Sub, opts *core.Options) error {
	t.client = client
	t.winevts = sub.WindowChanges()
	t.opts = opts
	return nil
}

func (t *Tiler) Configure(args []string) error {
	return nil
}

func (t *Tiler) Run() {
	for evt := range t.winevts {
		if evt.Change != ipc.FocusWindow {
			continue
		}

		if evt.Container.Type == ipc.FloatingConNode {
			if t.opts.Debug {
				t.opts.Log.Message("tiler", "Focused window is floating")
			}
			continue
		}

		t.opts.Log.Messagef("tiler", "Window: %v, Layout: %v", evt.Container.Name, evt.Container.Layout)

		root, err := t.client.Tree()
		if err != nil {
			t.opts.Log.Messagef("tiler", "Fatal error: GetTree failed: %v", err)
			break
		}

		parent := core.FindParent(root, evt.Container.ID)
		if parent == nil {
			if t.opts.Debug {
				t.opts.Log.Message("tiler", "Window has no parent")
			}
			continue
		}

		t.opts.Log.Messagef("tiler", "Parent: %v, Layout: %v", parent.Type, parent.Layout)

		if parent.Layout == "stacked" || parent.Layout == "tabbed" {
			if t.opts.Debug {
				t.opts.Log.Messagef("tiler", "Parent layout is excluded: %v", parent.Layout)
			}
			continue
		}

		newlayout := "splith"
		if evt.Container.Rect.Height > evt.Container.Rect.Width {
			newlayout = "splitv"
		}

		if t.opts.Debug {
			t.opts.Log.Messagef("tiler", "Selecting layout: %v", newlayout)
		}

		if parent.Layout != newlayout {
			s, err := t.client.Command(newlayout)
			if err != nil {
				t.opts.Log.Messagef("tiler", "Error sending command: %v", err)
			}
			if !(s[0].Success) {
				t.opts.Log.Messagef("tiler", "sway error: %v", s[0].Error)
			}
		}
	}

	if t.opts.Debug {
		t.opts.Log.Message("tiler", "Tiling channel closed")
	}
}

func (t *Tiler) Close() {
	t.opts.Log.Message("tiler", "Closing Block")
}
