/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package blockpalette

import (
	"strings"
	"sync"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

var (
	stairsOnce sync.Once
	stairsFwd  map[uint32]uint32 // dragonfly rid -> 1.26.50 hash
	stairsRev  map[uint32]uint32 // 1.26.50 hash -> dragonfly rid
)

func isNative(proto uint32) bool { return proto == uint32(protocol.CurrentProtocol) }

func buildStairs() {
	stairsFwd, stairsRev = map[uint32]uint32{}, map[uint32]uint32{}
	reg := world.DefaultBlockRegistry
	for rid := uint32(0); rid < uint32(reg.BlockCount()); rid++ {
		name, props, ok := reg.RuntimeIDToState(rid)
		if !ok || !strings.HasSuffix(name, "_stairs") {
			continue
		}
		if _, has := props["minecraft:corner"]; has {
			continue // already the new shape
		}
		np := map[string]any{"minecraft:corner": "none"} // added to stairs in 1.26.50
		for k, v := range props {
			np[k] = v
		}
		h := world.NetworkBlockHash(name, np)
		stairsFwd[rid], stairsRev[h] = h, rid
	}
}

func nativeStairHash(rid uint32) uint32 {
	stairsOnce.Do(buildStairs)
	return stairsFwd[rid]
}

func nativeStairRID(hash uint32) (uint32, bool) {
	stairsOnce.Do(buildStairs)
	rid, ok := stairsRev[hash]
	return rid, ok
}
