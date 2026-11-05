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

package v630

import (
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Recipe struct {
	Kind      int32
	Shaped    protocol.ShapedRecipe
	Shapeless protocol.ShapelessRecipe
	Multi     protocol.MultiRecipe
	Transform protocol.SmithingTransformRecipe
	Trim      protocol.SmithingTrimRecipe
}

func (x *Recipe) Marshal(io protocol.IO) {
	io.Varint32(&x.Kind)
	switch x.Kind {
	case 0, 5, 6:
		r := &x.Shapeless
		io.String(&r.RecipeID)
		protocol.FuncSlice(io, &r.Input, io.ItemDescriptorCount)
		protocol.FuncSlice(io, &r.Output, io.Item)
		io.UUID(&r.UUID)
		io.String(&r.Block)
		io.Varint32(&r.Priority)
		io.Varuint32(&r.RecipeNetworkID)
	case 1, 7:
		r := &x.Shaped
		io.String(&r.RecipeID)
		io.Varint32(&r.Width)
		io.Varint32(&r.Height)
		if r.Width < 1 || r.Width > 3 || r.Height < 1 || r.Height > 3 {
			io.InvalidValue(r.Width, "recipe shape", "dimensions outside crafting grid")
			return
		}
		protocol.FuncSliceOfLen(io, uint32(r.Width*r.Height), &r.Input, io.ItemDescriptorCount)
		protocol.FuncSlice(io, &r.Output, io.Item)
		io.UUID(&r.UUID)
		io.String(&r.Block)
		io.Varint32(&r.Priority)
		io.Varuint32(&r.RecipeNetworkID)
	case 4:
		protocol.Single(io, &x.Multi)
	case 8:
		protocol.Single(io, &x.Transform)
	case 9:
		protocol.Single(io, &x.Trim)
	default:
		io.UnknownEnumOption(x.Kind, "legacy recipe type")
	}
}

type CraftingData struct {
	Recipes                      []Recipe
	PotionRecipes                []protocol.PotionRecipe
	PotionContainerChangeRecipes []protocol.PotionContainerChangeRecipe
	MaterialReducers             []protocol.MaterialReducer
	ClearRecipes                 bool
}

func (*CraftingData) ID() uint32 { return packet.IDCraftingData }

func (pk *CraftingData) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.Recipes)
	protocol.Slice(io, &pk.PotionRecipes)
	protocol.Slice(io, &pk.PotionContainerChangeRecipes)
	protocol.FuncSlice(io, &pk.MaterialReducers, io.MaterialReducer)
	io.Bool(&pk.ClearRecipes)
}

func fromLatestCrafting(proto uint32, pk *packet.CraftingData) *CraftingData {
	out := &CraftingData{ClearRecipes: pk.ClearRecipes}
	validInputs := func(items []protocol.ItemDescriptorCount) bool {
		for _, in := range items {
			if d, ok := in.Descriptor.(*protocol.DefaultItemDescriptor); ok {
				native, ok := world.ItemRuntimeIDByName(d.Name)
				if !ok {
					return false
				}
				if id, _ := itemdata.TranslateStack(proto, native, uint32(d.MetadataValue)); id == 0 {
					return false
				}
			}
		}
		return true
	}
	validOutput := func(items []protocol.ItemStack) bool {
		for _, it := range items {
			if id, _ := itemdata.TranslateStack(proto, it.NetworkID, it.MetadataValue); id == 0 {
				return false
			}
		}
		return true
	}
	addShaped := func(kind int32, r protocol.ShapedRecipe) {
		if validInputs(r.Input) && validOutput(r.Output) {
			out.Recipes = append(out.Recipes, Recipe{Kind: kind, Shaped: r})
		}
	}
	addShapeless := func(kind int32, r protocol.ShapelessRecipe) {
		if validInputs(r.Input) && validOutput(r.Output) {
			out.Recipes = append(out.Recipes, Recipe{Kind: kind, Shapeless: r})
		}
	}
	for _, r := range pk.ShapedRecipes {
		addShaped(1, r)
	}
	for _, r := range pk.ShapelessRecipes {
		addShapeless(0, r)
	}
	for _, r := range pk.MultiRecipes {
		out.Recipes = append(out.Recipes, Recipe{Kind: 4, Multi: r})
	}
	for _, r := range pk.UserDataShapelessRecipes {
		addShapeless(5, r.ShapelessRecipe)
	}
	for _, r := range pk.ShapelessChemistryRecipes {
		addShapeless(6, r.ShapelessRecipe)
	}
	for _, r := range pk.ShapedChemistryRecipes {
		addShaped(7, r.ShapedRecipe)
	}
	for _, r := range pk.SmithingTransformRecipes {
		if validInputs([]protocol.ItemDescriptorCount{r.Template, r.Base, r.Addition}) && validOutput([]protocol.ItemStack{r.Result}) {
			out.Recipes = append(out.Recipes, Recipe{Kind: 8, Transform: r})
		}
	}
	for _, r := range pk.SmithingTrimRecipes {
		if validInputs([]protocol.ItemDescriptorCount{r.Template, r.Base, r.Addition}) {
			out.Recipes = append(out.Recipes, Recipe{Kind: 9, Trim: r})
		}
	}
	for _, r := range pk.PotionRecipes {
		input, inputMeta := itemdata.TranslateStack(proto, r.InputPotionID, uint32(r.InputPotionMetadata))
		reagent, reagentMeta := itemdata.TranslateStack(proto, r.ReagentItemID, uint32(r.ReagentItemMetadata))
		output, outputMeta := itemdata.TranslateStack(proto, r.OutputPotionID, uint32(r.OutputPotionMetadata))
		if input != 0 && reagent != 0 && output != 0 {
			out.PotionRecipes = append(out.PotionRecipes, protocol.PotionRecipe{InputPotionID: input, InputPotionMetadata: int32(inputMeta), ReagentItemID: reagent, ReagentItemMetadata: int32(reagentMeta), OutputPotionID: output, OutputPotionMetadata: int32(outputMeta)})
		}
	}
	for _, r := range pk.PotionContainerChangeRecipes {
		input, _ := itemdata.TranslateStack(proto, r.InputItemID, 0)
		reagent, _ := itemdata.TranslateStack(proto, r.ReagentItemID, 0)
		output, _ := itemdata.TranslateStack(proto, r.OutputItemID, 0)
		if input != 0 && reagent != 0 && output != 0 {
			out.PotionContainerChangeRecipes = append(out.PotionContainerChangeRecipes, protocol.PotionContainerChangeRecipe{InputItemID: input, ReagentItemID: reagent, OutputItemID: output})
		}
	}
	for _, r := range pk.MaterialReducers {
		input, meta := itemdata.TranslateStack(proto, r.InputItem.NetworkID, r.InputItem.MetadataValue)
		if input == 0 {
			continue
		}
		mapped := protocol.MaterialReducer{InputItem: protocol.ItemType{NetworkID: input, MetadataValue: meta}}
		for _, result := range r.Outputs {
			id, _ := itemdata.TranslateStack(proto, result.NetworkID, 0)
			if id == 0 {
				break
			}
			mapped.Outputs = append(mapped.Outputs, protocol.MaterialReducerOutput{NetworkID: id, Count: result.Count})
		}
		if len(mapped.Outputs) == len(r.Outputs) {
			out.MaterialReducers = append(out.MaterialReducers, mapped)
		}
	}
	return out
}
