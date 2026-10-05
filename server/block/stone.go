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

package block

import "github.com/df-mc/dragonfly/server/item"

type (
	Stone struct {
		solid
		bassDrum

		Smooth bool
	}

	Granite polishable

	Diorite polishable

	Andesite polishable

	polishable struct {
		solid
		bassDrum

		Polished bool
	}
)

func (s Stone) BreakInfo() BreakInfo {
	if s.Smooth {
		return newBreakInfo(2, pickaxeHarvestable, pickaxeEffective, oneOf(s)).withBlastResistance(6)
	}
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, silkTouchOneOf(Cobblestone{}, Stone{})).withBlastResistance(6)
}

func (g Granite) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, oneOf(g)).withBlastResistance(6)
}

func (d Diorite) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, oneOf(d)).withBlastResistance(6)
}

func (a Andesite) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, oneOf(a)).withBlastResistance(6)
}

func (s Stone) SmeltInfo() item.SmeltInfo {
	if s.Smooth {
		return item.SmeltInfo{}
	}
	return newSmeltInfo(item.NewStack(Stone{Smooth: true}, 1), 0.1)
}

func (s Stone) EncodeItem() (name string, meta int16) {
	if s.Smooth {
		return "minecraft:smooth_stone", 0
	}
	return "minecraft:stone", 0
}

func (s Stone) EncodeBlock() (string, map[string]any) {
	if s.Smooth {
		return "minecraft:smooth_stone", nil
	}
	return "minecraft:stone", nil
}

func (a Andesite) EncodeItem() (name string, meta int16) {
	if a.Polished {
		return "minecraft:polished_andesite", 0
	}
	return "minecraft:andesite", 0
}

func (a Andesite) EncodeBlock() (string, map[string]any) {
	if a.Polished {
		return "minecraft:polished_andesite", nil
	}
	return "minecraft:andesite", nil
}

func (d Diorite) EncodeItem() (name string, meta int16) {
	if d.Polished {
		return "minecraft:polished_diorite", 0
	}
	return "minecraft:diorite", 0
}

func (d Diorite) EncodeBlock() (string, map[string]any) {
	if d.Polished {
		return "minecraft:polished_diorite", nil
	}
	return "minecraft:diorite", nil
}

func (g Granite) EncodeItem() (name string, meta int16) {
	if g.Polished {
		return "minecraft:polished_granite", 0
	}
	return "minecraft:granite", 0
}

func (g Granite) EncodeBlock() (string, map[string]any) {
	if g.Polished {
		return "minecraft:polished_granite", nil
	}
	return "minecraft:granite", nil
}
