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

package blockinternal

import (
	"github.com/df-mc/dragonfly/server/item/category"
	"maps"
	"slices"
)

type ComponentBuilder struct {
	permutations map[string]map[string]any
	properties   []map[string]any
	components   map[string]any
	blockID      int32

	identifier   string
	menuCategory category.Category
}

func NewComponentBuilder(identifier string, components map[string]any, blockID int32) *ComponentBuilder {
	if components == nil {
		components = map[string]any{}
	}
	return &ComponentBuilder{
		permutations: make(map[string]map[string]any),
		components:   components,
		blockID:      blockID,

		identifier:   identifier,
		menuCategory: category.Construction(),
	}
}

func (builder *ComponentBuilder) AddProperty(name string, values []any) {
	builder.properties = append(builder.properties, map[string]any{
		"name": name,
		"enum": values,
	})
}

func (builder *ComponentBuilder) AddComponent(name string, value any) {
	builder.components[name] = value
}

func (builder *ComponentBuilder) AddPermutation(condition string, components map[string]any) {
	if len(builder.permutations) == 0 {

		builder.AddComponent("minecraft:on_player_placing", map[string]any{
			"triggerType": "placement_trigger",
		})
	}
	if builder.permutations[condition] == nil {
		builder.permutations[condition] = map[string]any{}
	}
	for key, value := range components {
		builder.permutations[condition][key] = value
	}
}

func (builder *ComponentBuilder) SetMenuCategory(category category.Category) {
	builder.menuCategory = category
}

func (builder *ComponentBuilder) Construct() map[string]any {
	properties := slices.Clone(builder.properties)
	components := maps.Clone(builder.components)

	result := map[string]any{
		"components":    components,
		"molangVersion": int32(10),
		"menu_category": map[string]any{
			"category": builder.menuCategory.String(),
			"group":    builder.menuCategory.Group(),
		},
		"vanilla_block_data": map[string]any{
			"block_id": builder.blockID,
		},
	}
	if len(properties) > 0 {
		result["properties"] = properties
	}

	permutations := maps.Clone(builder.permutations)
	if len(permutations) > 0 {
		result["permutations"] = []map[string]any{}
		for condition, values := range permutations {
			result["permutations"] = append(result["permutations"].([]map[string]any), map[string]any{
				"condition":  condition,
				"components": values,
			})
		}
	}
	return result
}
