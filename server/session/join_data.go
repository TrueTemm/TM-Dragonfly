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

package session

import (
	"sync"

	"github.com/df-mc/dragonfly/server/item/recipe"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type joinData struct {
	creative *packet.CreativeContent
	crafting *packet.CraftingData
	recipes  map[uint32]recipe.Recipe
	biomes   *packet.BiomeDefinitionList
	trim     *packet.TrimData
}

var (
	joinDataMu    sync.Mutex
	joinDataCache = map[world.BlockRegistry]*joinData{}
)

func joinDataFor(br world.BlockRegistry) *joinData {
	joinDataMu.Lock()
	defer joinDataMu.Unlock()
	if d, ok := joinDataCache[br]; ok && len(d.recipes) == len(recipe.Recipes()) {
		return d
	}
	groups, items := creativeContent(br)
	crafting, recipes := buildCraftingData(br)
	d := &joinData{
		creative: &packet.CreativeContent{Groups: groups, Items: items},
		crafting: crafting,
		recipes:  recipes,
		biomes:   buildBiomeDefinitionList(),
		trim:     buildTrimData(),
	}
	joinDataCache[br] = d
	return d
}
