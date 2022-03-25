// adapted from https://github.com/nwg-piotr/autotiling/blob/master/autotiling/main.py
package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type Tiler struct {
	client   core.Client
	winevts  ipc.Cookie
	modevts  ipc.Cookie
	opts     *core.Options
	loglevel core.LogLevel
}

func init() {
	var _ core.BlockInitializer = (*Tiler)(nil)
}

func (t *Tiler) Init(client core.Client, sub core.Sub, opts *core.Options, args ...string) error {
	t.client = client
	wincookie, err := sub.WindowChanges(t)
	if err != nil {
		return err
	}

	modcookie, err := sub.BindingChanges(t)
	if err != nil {
		return err
	}

	t.opts = opts
	t.winevts = wincookie
	t.modevts = modcookie
	return nil
}

func (t *Tiler) SetLogLevel(level core.LogLevel) {
	t.loglevel = level
}

func (t *Tiler) BindingChange(evt ipc.BindingChange) {
	setLayout(t)
}

func (t *Tiler) WindowChange(evt ipc.WindowChange) {
	if evt.Change != ipc.FocusWindow {
		return
	}

	setLayout(t)
}

func setLayout(t *Tiler) {
	root, err := t.client.Tree()
	if err != nil {
		t.opts.Log.Printf("tiler", "Fatal error: GetTree failed: %v", err)
		return
	}

	con := root.FindChild(func(n *ipc.Node) bool {
		return n.Focused
	})

	if con == nil {
		if t.loglevel.Debug() {
			t.opts.Log.Print("tiler", "Focused window not found")
		}
		return
	}

	if con.Type == ipc.FloatingConNode ||
		con.FullScreenMode != ipc.FullScreenModeNone {
		return
	}

	parent := core.FindParent(root, con.ID)
	if parent == nil {
		if t.loglevel.Debug() {
			t.opts.Log.Print("tiler", "Window has no parent")
		}
		return
	}

	if parent.Layout == "stacked" || parent.Layout == "tabbed" {
		return
	}

	newlayout := "splith"
	if con.Rect.Height > con.Rect.Width {
		newlayout = "splitv"
	}

	if parent.Layout != newlayout {
		s, err := t.client.Command(newlayout)
		if err != nil {
			t.opts.Log.Printf("tiler", "Error sending command: %v", err)
		}

		if len(s) > 0 {
			if s[0].Success {
				if t.loglevel.Debug() {
					t.opts.Log.Printf("tiller", "switched to layout: %v", newlayout)
				}
			} else {
				t.opts.Log.Printf("tiler", "sway error: %v", s[0].Error)
			}
		}
	}
	return
}
