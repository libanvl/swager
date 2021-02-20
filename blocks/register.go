package blocks

import "github.com/libanvl/swager/internal/core"

func RegisterBlocks() {
	core.Blocks.Register("swaymon",
		func() core.BlockInitializer { return new(SwayMon) })
	core.Blocks.Register("tiler",
		func() core.BlockInitializer { return new(Tiler) })
	core.Blocks.Register("initspawn",
		func() core.BlockInitializer { return new(InitSpawn) })
	core.Blocks.Register("execnew",
		func() core.BlockInitializer { return new(ExecNew) })
}
