package blocks

import (
	"log"

	"github.com/libanvl/swager/internal/core"
	"go.i3wm.org/i3/v4"
)

type Tiler struct {
}

func init() {
	var _ core.ChangeEventBlock = (*Tiler)(nil)
}

func (t *Tiler) Init(wm core.WinMgr) error {
  log.Printf("Initializing Tiler Block for winmgr: %s", wm.String())
	return nil
}

func (t *Tiler) Configure(args []string) error {
	log.Print("Configuring Tiler Block")
	return nil
}

func (t *Tiler) Event() []i3.EventType {
	return []i3.EventType{i3.WindowEventType}
}

func (t *Tiler) MatchChange(change string) bool {
	return change == "focus"
}

func (t *Tiler) OnEvent(evt interface{}) error {
	windowevt := evt.(*i3.WindowEvent)
	log.Printf("Window: %v, Layout: %v", windowevt.Container.Name, windowevt.Container.Layout)
	parent := core.FindParent(windowevt.Container.ID)
	log.Printf("Parent: %v, Layout: %v", parent.Type, parent.Layout)
	log.Printf("%#v", windowevt)
  return nil
}

func (t *Tiler) Close() {
	log.Print("Closing Tiler Block")
}
