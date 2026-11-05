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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestLegacyCraftingWire(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	stone, _ := world.ItemRuntimeIDByName("minecraft:stone")
	output := []protocol.ItemStack{{ItemType: protocol.ItemType{NetworkID: stone}, Count: 1}}
	source := &packet.CraftingData{ClearRecipes: true,
		ShapedRecipes:    []protocol.ShapedRecipe{{RecipeID: "test_shaped", Width: 1, Height: 1, Input: []protocol.ItemDescriptorCount{{Descriptor: &protocol.DefaultItemDescriptor{Name: "minecraft:oak_planks"}, Count: 1}}, Output: output, Block: "crafting_table", RecipeNetworkID: 11}},
		ShapelessRecipes: []protocol.ShapelessRecipe{{RecipeID: "test_shapeless", Input: []protocol.ItemDescriptorCount{{Descriptor: &protocol.ItemTagItemDescriptor{Tag: "minecraft:planks"}, Count: 1}}, Output: output, Block: "crafting_table", RecipeNetworkID: 12}},
	}
	for _, id := range []int32{589, 594, 618, 622, 630} {
		t.Run(fmt.Sprint(id), func(t *testing.T) {
			data, err := os.ReadFile(fmt.Sprintf("testdata/wire%d.json", id))
			if err != nil {
				t.Fatal(err)
			}
			var fixtures map[string][]byte
			if err := json.Unmarshal(data, &fixtures); err != nil {
				t.Fatal(err)
			}
			p := Protocol{Version: id}
			converted := p.ConvertFromLatest(source, nil)[0].(*CraftingData)
			if len(converted.Recipes) != 2 {
				t.Fatal("recipes were discarded")
			}
			buf := new(bytes.Buffer)
			converted.Marshal(p.NewWriter(buf, 0))
			if !bytes.Equal(buf.Bytes(), fixtures["crafting"]) {
				t.Fatalf("historical crafting wire differs: got %x want %x", buf.Bytes(), fixtures["crafting"])
			}
			var decoded CraftingData
			decoded.Marshal(p.NewReader(buf, 0, true))
			ingredient := decoded.Recipes[0].Shaped.Input[0].Descriptor.(*protocol.DefaultItemDescriptor)
			if len(decoded.Recipes) != 2 || ingredient.Name != "minecraft:oak_planks" || buf.Len() != 0 {
				t.Fatal("crafting ingredients did not survive decoding")
			}
		})
	}
	if source.ShapedRecipes[0].Input[0].Descriptor.(*protocol.DefaultItemDescriptor).Name != "minecraft:oak_planks" {
		t.Fatal("source recipe changed")
	}
}

func TestLegacyMaterialReducerWire(t *testing.T) {
	stone, _ := world.ItemRuntimeIDByName("minecraft:stone")
	coarse, _ := world.ItemRuntimeIDByName("minecraft:coarse_dirt")
	unknown, _ := world.ItemRuntimeIDByName("minecraft:mace")
	for _, proto := range []int32{589, 594, 618, 622, 630} {
		data, err := os.ReadFile(fmt.Sprintf("testdata/wire%d.json", proto))
		if err != nil {
			t.Fatal(err)
		}
		var fixtures map[string][]byte
		if err := json.Unmarshal(data, &fixtures); err != nil {
			t.Fatal(err)
		}
		valid := protocol.MaterialReducer{InputItem: protocol.ItemType{NetworkID: coarse}, Outputs: []protocol.MaterialReducerOutput{{NetworkID: stone, Count: 2}}}
		source := &packet.CraftingData{ClearRecipes: true, MaterialReducers: []protocol.MaterialReducer{valid,
			{InputItem: protocol.ItemType{NetworkID: unknown}, Outputs: valid.Outputs},
			{InputItem: valid.InputItem, Outputs: []protocol.MaterialReducerOutput{{NetworkID: unknown, Count: 1}}},
		}}
		p := Protocol{Version: proto}
		out := p.ConvertFromLatest(source, nil)[0].(*CraftingData)
		if len(out.MaterialReducers) != 1 || source.MaterialReducers[0].InputItem.NetworkID != coarse {
			t.Fatal("unsupported reducer or source mutation")
		}
		buf := new(bytes.Buffer)
		out.Marshal(p.NewWriter(buf, 0))
		if !bytes.Equal(buf.Bytes(), fixtures["material_reducers"]) {
			t.Fatalf("protocol %d reducer differs from historical codec", proto)
		}
	}
}
