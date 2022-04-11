package blocks

import (
	"encoding/json"

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
	var _ ipc.BindingModeChangeHandler = (*Tiler)(nil)
	var _ ipc.WindowChangeHandler = (*Tiler)(nil)
}

func (t *Tiler) Init(client core.Client, sub core.Sub, opts *core.Options, args ...string) error {
	t.client = client
	wincookie, err := sub.WindowChanges(t)
	if err != nil {
		return err
	}

	modcookie, err := sub.BindingModeChanges(t)
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

func (t *Tiler) BindingModeChange(evt ipc.ModeChange) {
	if t.loglevel.Debug() {
		t.opts.Log.Printf("tiler", "Binding Mode Change event: %v", evt)
	}
	setLayout(t)
}

func (t *Tiler) WindowChange(evt ipc.WindowChange) {
	if t.loglevel.Debug() {
		t.opts.Log.Printf("tiler", "Window Change event: %v %v", evt.Change, evt.Container.Name)
	}

	switch evt.Change {
	case ipc.TitleWindow:
	case ipc.MarkWindow:
	case ipc.UrgentWindow:
		return
	}

	setLayout(t)
}

//adapted from https://github.com/nwg-piotr/autotiling/blob/master/autotiling/main.py
func setLayout(t *Tiler) {
	root, err := t.client.Tree()
	if err != nil {
		t.opts.Log.Printf("tiler", "GetTree failed: %v", err)
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
		*con.FullscreenMode != ipc.NoneFullscreenMode {
		return
	}

	parent := core.FindParent(root, con.ID)
	if parent == nil {
		if t.loglevel.Debug() {
			t.opts.Log.Print("tiler", "Window has no parent")
		}
		return
	}

	if parent.Layout == ipc.StackedLayout || parent.Layout == ipc.TabbedLayout {
		return
	}

	newlayout := ipc.SplitHLayout
	if con.Rect.Height > con.Rect.Width {
		newlayout = ipc.SplitVLayout
	}

	if parent.Layout != newlayout {
		s, err := t.client.Command(string(newlayout))
		if err != nil {
			if jerr, ok := err.(*json.UnmarshalTypeError); ok {
				t.opts.Log.Printf("tiler", "Error sending command: %v", jerr)
			}
			t.opts.Log.Printf("tiler", "Error sending command: %v", err)
		}

		if len(s) > 0 {
			if s[0].Success {
				if t.loglevel.Debug() {
					t.opts.Log.Printf("tiler", "switched to layout: %v", newlayout)
				}
			} else {
				t.opts.Log.Printf("tiler", "sway error: %v", s[0].Error)
			}
		}
	}
	return
}
