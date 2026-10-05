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
	"github.com/df-mc/dragonfly/server/internal/sliceutil"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"slices"
	"sort"
	"strings"
	"unsafe"
)

var (
	recipes []Recipe

	dynamicRecipes []DynamicRecipe

	index = make(map[string]map[string]Recipe)

	reagent = make(map[string]item.Stack)
)

func Recipes() []Recipe {
	return slices.Clone(recipes)
}

func DynamicRecipes() []DynamicRecipe {
	return slices.Clone(dynamicRecipes)
}

func Register(recipe Recipe) {
	recipes = append(recipes, recipe)

	_, ok := recipe.(PotionContainerChange)
	p, okTwo := recipe.(Potion)

	if okTwo {
		stack := p.Input()[1].(item.Stack)
		name, _ := stack.Item().EncodeItem()
		reagent[name] = stack
	}

	if ok || okTwo {
		input := make([]world.Item, len(recipe.Input()))
		for i, stack := range recipe.Input() {
			if s, ok := stack.(item.Stack); ok {
				input[i] = s.Item()
			}
		}
		hash := hashItems(input, !ok)

		block := recipe.Block()
		if index[block] == nil {
			index[block] = make(map[string]Recipe)
		}
		index[block][hash] = recipe
	}
}

func Perform(block string, input ...world.Item) (output []item.Stack, ok bool) {
	blockInd, ok := index[block]
	if !ok {

		return nil, false
	}
	r, ok := blockInd[hashItems(input, true)]
	if !ok {
		r, ok = blockInd[hashItems(input, false)]
		if !ok {
			return nil, false
		}
	}
	_, containerChange := r.(PotionContainerChange)
	for ind, it := range r.Output() {
		if containerChange {
			name, _ := it.Item().EncodeItem()
			_, meta := input[ind].EncodeItem()
			if i, ok := world.ItemByName(name, meta); ok {
				it = item.NewStack(i, it.Count())
			}
		}
		output = append(output, it)
	}
	return output, ok
}

func hashItems(items []world.Item, useMeta bool) string {
	items = sliceutil.Filter(items, func(it world.Item) bool {
		return it != nil
	})
	sort.Slice(items, func(i, j int) bool {
		nameOne, metaOne := items[i].EncodeItem()
		nameTwo, metaTwo := items[j].EncodeItem()
		if nameOne == nameTwo {
			return metaOne < metaTwo
		}
		return nameOne < nameTwo
	})

	var b strings.Builder
	for _, it := range items {
		name, meta := it.EncodeItem()
		b.WriteString(name)
		if useMeta {
			a := *(*[2]byte)(unsafe.Pointer(&meta))
			b.Write(a[:])
		}
	}
	return b.String()
}

func ValidBrewingReagent(i world.Item) bool {
	name, _ := i.EncodeItem()
	_, exists := reagent[name]
	return exists
}

func RegisterDynamic(recipe DynamicRecipe) {
	dynamicRecipes = append(dynamicRecipes, recipe)
}
