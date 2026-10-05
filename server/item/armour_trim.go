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

package item

import "github.com/df-mc/dragonfly/server/world"

type ArmourTrim struct {
	Template SmithingTemplateType
	Material ArmourTrimMaterial
}

func (trim ArmourTrim) Zero() bool {
	return trim.Material == nil || trim.Template == TemplateNetheriteUpgrade()
}

type ArmourTrimMaterial interface {
	TrimMaterial() string

	MaterialColour() string
}

func trimMaterialFromString(name string) (ArmourTrimMaterial, bool) {
	switch name {
	case "amethyst":
		return AmethystShard{}, true
	case "copper":
		return CopperIngot{}, true
	case "diamond":
		return Diamond{}, true
	case "emerald":
		return Emerald{}, true
	case "gold":
		return GoldIngot{}, true
	case "iron":
		return IronIngot{}, true
	case "lapis":
		return LapisLazuli{}, true
	case "netherite":
		return NetheriteIngot{}, true
	case "quartz":
		return NetherQuartz{}, true
	case "resin":
		return ResinBrick{}, true
	case "redstone":
		return RedstoneWire{}, true
	}
	return nil, false
}

func ArmourTrimMaterials() []world.Item {
	return []world.Item{
		AmethystShard{},
		CopperIngot{},
		Diamond{},
		Emerald{},
		GoldIngot{},
		IronIngot{},
		LapisLazuli{},
		NetheriteIngot{},
		NetherQuartz{},
		ResinBrick{},
		RedstoneWire{},
	}
}

type Trimmable interface {
	WithTrim(trim ArmourTrim) world.Item
}
