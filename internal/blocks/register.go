package blocks

import "github.com/libanvl/swager/internal/core"

func RegisterBlocks() {
	core.Blocks.Register("tiler",
		func() core.Block { return new(Tiler) })
	core.Blocks.Register("initspawn",
		func() core.Block { return new(InitSpawn) })
}
