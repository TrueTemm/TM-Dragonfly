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
	"github.com/google/uuid"
)

type Recipe interface {
	Input() []Item

	Output() []item.Stack

	Block() string

	Priority() uint32
}

type DynamicRecipe interface {
	Match(input []Item) (output []item.Stack, ok bool)

	Block() string
}

type Shapeless struct {
	recipe
}

func NewShapeless(input []Item, output item.Stack, block string) Shapeless {
	return Shapeless{recipe: recipe{
		input:  input,
		output: []item.Stack{output},
		block:  block,
	}}
}

type UserDataShapeless struct {
	recipe
}

func NewUserDataShapeless(input []Item, output item.Stack, block string) UserDataShapeless {
	return UserDataShapeless{recipe: recipe{
		input:  input,
		output: []item.Stack{output},
		block:  block,
	}}
}

type Multi struct {
	recipe
	uuid uuid.UUID
}

func NewMulti(id uuid.UUID) Multi {
	return Multi{uuid: id, recipe: recipe{block: "crafting_table"}}
}

func (m Multi) UUID() uuid.UUID {
	return m.uuid
}

type SmithingTransform struct {
	recipe
}

func NewSmithingTransform(base, addition, template Item, output item.Stack, block string) SmithingTransform {
	return SmithingTransform{recipe: recipe{
		input:  []Item{base, addition, template},
		output: []item.Stack{output},
		block:  block,
	}}
}

type SmithingTrim struct {
	recipe
}

func NewSmithingTrim(base, addition, template Item, block string) SmithingTrim {
	return SmithingTrim{recipe: recipe{
		input: []Item{base, addition, template},
		block: block,
	}}
}

type PotionContainerChange struct {
	recipe
}

func NewPotionContainerChange(input, output world.Item, reagent item.Stack) PotionContainerChange {
	return PotionContainerChange{recipe: recipe{
		input:  []Item{item.NewStack(input, 1), reagent},
		output: []item.Stack{item.NewStack(output, 1)},
		block:  "brewing_stand",
	}}
}

type Potion struct {
	recipe
}

func NewPotion(input, reagent Item, output item.Stack) Potion {
	return Potion{recipe: recipe{
		input:  []Item{input, reagent},
		output: []item.Stack{output},
		block:  "brewing_stand",
	}}
}

type Shaped struct {
	recipe

	shape Shape
}

func NewShaped(input []Item, output item.Stack, shape Shape, block string) Shaped {
	return Shaped{
		shape: shape,
		recipe: recipe{
			input:  input,
			output: []item.Stack{output},
			block:  block,
		},
	}
}

func (r Shaped) Shape() Shape {
	return r.shape
}

type recipe struct {
	input []Item

	output []item.Stack

	block string

	priority uint32
}

func (r recipe) Input() []Item {
	return r.input
}

func (r recipe) Output() []item.Stack {
	return r.output
}

func (r recipe) Block() string {
	return r.block
}

func (r recipe) Priority() uint32 {
	return r.priority
}
