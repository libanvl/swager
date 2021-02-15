package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type Tiler struct {
	client   core.Client
	winevts  <-chan *ipc.WindowChange
	opts     *core.Options
	loglevel core.LogLevel
}

func init() {
	var _ core.Block = (*Tiler)(nil)
}

func (t *Tiler) Init(client core.Client, sub core.Sub, opts *core.Options, args ...string) error {
	t.client = client
	t.winevts = sub.WindowChanges()
	t.opts = opts
	return nil
}

func (t *Tiler) SetLogLevel(level core.LogLevel) {
	t.loglevel = level
}

func (t *Tiler) Run() {
	for evt := range t.winevts {
		if evt.Change != ipc.FocusWindow {
			continue
		}

		if evt.Container.Type == ipc.FloatingConNode {
			continue
		}

		root, err := t.client.Tree()
		if err != nil {
			t.opts.Log.Printf("tiler", "Fatal error: GetTree failed: %v", err)
			break
		}

		parent := core.FindParent(root, evt.Container.ID)
		if parent == nil {
			if t.loglevel.Debug() {
				t.opts.Log.Print("tiler", "Window has no parent")
			}
			continue
		}

		if parent.Layout == "stacked" || parent.Layout == "tabbed" {
			continue
		}

		newlayout := "splith"
		if evt.Container.Rect.Height > evt.Container.Rect.Width {
			newlayout = "splitv"
		}

		if parent.Layout != newlayout {
			s, err := t.client.Command(newlayout)
			if err != nil {
				t.opts.Log.Printf("tiler", "Error sending command: %v", err)
			}
			if !(s[0].Success) {
				t.opts.Log.Printf("tiler", "sway error: %v", s[0].Error)
			}
		}
	}
}

func (t *Tiler) Close() {
}
