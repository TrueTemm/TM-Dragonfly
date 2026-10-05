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

package recipe

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type DecoratedPotRecipe struct {
	block string
}

func NewDecoratedPotRecipe() DecoratedPotRecipe {
	return DecoratedPotRecipe{block: "crafting_table"}
}

type potDecoration interface {
	world.Item
	PotDecoration() bool
}

func (r DecoratedPotRecipe) Match(input []Item) (output []item.Stack, ok bool) {

	if len(input) != 9 {
		return nil, false
	}

	decorations := [4]world.Item{}
	decorationIndex := 0
	for i := range input {
		it := input[i]
		if i%2 == 0 {

			if !it.Empty() {
				return nil, false
			}
		} else {

			if it.Empty() {
				return nil, false
			}

			var actualItem item.Stack
			if v, ok := it.(item.Stack); ok {
				actualItem = v
			} else {

				return nil, false
			}

			decoration, ok := actualItem.Item().(potDecoration)
			if !ok {
				return nil, false
			}
			decorations[decorationIndex] = decoration
			decorationIndex++
		}
	}

	pot, ok := world.BlockByName("minecraft:decorated_pot", map[string]any{"direction": int32(2)})
	if !ok {
		return nil, false
	}

	sherds := []any{}

	for _, idx := range []int{0, 1, 3, 2} {
		name, _ := decorations[idx].EncodeItem()
		sherds = append(sherds, name)
	}

	if nbtDecoder, ok := pot.(interface {
		DecodeNBT(map[string]any) any
	}); ok {
		decodedPot := nbtDecoder.DecodeNBT(map[string]any{
			"id":     "DecoratedPot",
			"sherds": sherds,
		})
		if potItem, ok := decodedPot.(world.Item); ok {
			return []item.Stack{item.NewStack(potItem, 1)}, true
		}
	}

	return nil, false
}

func (r DecoratedPotRecipe) Block() string {
	return r.block
}
