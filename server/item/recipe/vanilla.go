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
	_ "embed"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

var (
	//go:embed crafting_data.nbt
	vanillaCraftingData []byte
	//go:embed smithing_data.nbt
	vanillaSmithingData []byte
	//go:embed smithing_trim_data.nbt
	vanillaSmithingTrimData []byte
	//go:embed potion_data.nbt
	vanillaPotionData []byte
)

type shapedRecipe struct {
	Input    inputItems  `nbt:"input"`
	Output   outputItems `nbt:"output"`
	Block    string      `nbt:"block"`
	Width    int32       `nbt:"width"`
	Height   int32       `nbt:"height"`
	Priority int32       `nbt:"priority"`
}

type shapelessRecipe struct {
	Input    inputItems  `nbt:"input"`
	Output   outputItems `nbt:"output"`
	Block    string      `nbt:"block"`
	Priority int32       `nbt:"priority"`
}

type potionRecipe struct {
	Input   inputItem  `nbt:"input"`
	Reagent inputItem  `nbt:"reagent"`
	Output  outputItem `nbt:"output"`
}

type potionContainerChangeRecipe struct {
	Input   string    `nbt:"input"`
	Reagent inputItem `nbt:"reagent"`
	Output  string    `nbt:"output"`
}

func registerVanilla() {
	var craftingRecipes struct {
		Shaped            []shapedRecipe    `nbt:"shaped"`
		Shapeless         []shapelessRecipe `nbt:"shapeless"`
		UserDataShapeless []shapelessRecipe `nbt:"shulker_box"`
		Multi             []string          `nbt:"multi"`
	}
	if err := nbt.Unmarshal(vanillaCraftingData, &craftingRecipes); err != nil {
		panic(err)
	}

	for _, id := range craftingRecipes.Multi {
		u, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		Register(NewMulti(u))
	}

	for _, s := range craftingRecipes.UserDataShapeless {
		input, ok := s.Input.Items()
		output, okTwo := s.Output.Stacks()
		if !ok || !okTwo {

			continue
		}
		Register(UserDataShapeless{recipe{
			input:    input,
			output:   output,
			block:    s.Block,
			priority: uint32(s.Priority),
		}})
	}

	for _, s := range craftingRecipes.Shapeless {
		input, ok := s.Input.Items()
		output, okTwo := s.Output.Stacks()
		if !ok || !okTwo {

			continue
		}
		Register(Shapeless{recipe{
			input:    input,
			output:   output,
			block:    s.Block,
			priority: uint32(s.Priority),
		}})
	}

	for _, s := range craftingRecipes.Shaped {
		input, ok := s.Input.Items()
		output, okTwo := s.Output.Stacks()
		if !ok || !okTwo {

			continue
		}
		Register(Shaped{
			shape: Shape{int(s.Width), int(s.Height)},
			recipe: recipe{
				input:    input,
				output:   output,
				block:    s.Block,
				priority: uint32(s.Priority),
			},
		})
	}

	var smithingRecipes []shapelessRecipe
	if err := nbt.Unmarshal(vanillaSmithingData, &smithingRecipes); err != nil {
		panic(err)
	}

	for _, s := range smithingRecipes {
		input, ok := s.Input.Items()
		output, okTwo := s.Output.Stacks()
		if !ok || !okTwo {

			continue
		}
		Register(SmithingTransform{recipe{
			input:    input,
			output:   output,
			block:    s.Block,
			priority: uint32(s.Priority),
		}})
	}

	var smithingTrimRecipes []shapelessRecipe
	if err := nbt.Unmarshal(vanillaSmithingTrimData, &smithingTrimRecipes); err != nil {
		panic(err)
	}

	for _, s := range smithingTrimRecipes {
		input, ok := s.Input.Items()
		if !ok {

			continue
		}
		Register(SmithingTrim{recipe{
			input:    input,
			block:    s.Block,
			priority: uint32(s.Priority),
		}})
	}

	var potionRecipes struct {
		Potions          []potionRecipe                `nbt:"potions"`
		ContainerChanges []potionContainerChangeRecipe `nbt:"container_changes"`
	}

	if err := nbt.Unmarshal(vanillaPotionData, &potionRecipes); err != nil {
		panic(err)
	}

	for _, r := range potionRecipes.Potions {
		input, ok := r.Input.Item()
		reagent, okTwo := r.Reagent.Item()
		output, okThree := r.Output.Stack()
		if !ok || !okTwo || !okThree {

			continue
		}

		Register(Potion{recipe{
			input:  []Item{input, reagent},
			output: []item.Stack{output},
			block:  "brewing_stand",
		}})
	}

	for _, c := range potionRecipes.ContainerChanges {
		input, ok := world.ItemByName(c.Input, 0)
		reagent, okTwo := c.Reagent.Item()
		output, okThree := world.ItemByName(c.Output, 0)
		if !ok || !okTwo || !okThree {

			continue
		}

		Register(PotionContainerChange{recipe{
			input:  []Item{item.NewStack(input, 1), reagent},
			output: []item.Stack{item.NewStack(output, 1)},
			block:  "brewing_stand",
		}})
	}

	RegisterDynamic(NewDecoratedPotRecipe())
}
