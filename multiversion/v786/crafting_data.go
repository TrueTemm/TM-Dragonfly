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

package v786

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type CraftingData struct {
	PotionRecipes []protocol.PotionRecipe

	PotionContainerChangeRecipes []protocol.PotionContainerChangeRecipe

	MaterialReducers []protocol.MaterialReducer

	ClearRecipes bool
}

func (*CraftingData) ID() uint32 {
	return IDCraftingData
}

func (pk *CraftingData) Marshal(io protocol.IO) {
	var noRecipes uint32
	io.Varuint32(&noRecipes)
	protocol.Slice(io, &pk.PotionRecipes)
	protocol.Slice(io, &pk.PotionContainerChangeRecipes)
	protocol.FuncSlice(io, &pk.MaterialReducers, io.MaterialReducer)
	io.Bool(&pk.ClearRecipes)
}
