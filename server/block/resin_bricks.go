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

type ResinBricks struct {
	solid
	bassDrum

	Chiseled bool
}

func (r ResinBricks) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, oneOf(r)).withBlastResistance(6)
}

func (r ResinBricks) EncodeItem() (name string, meta int16) {
	if r.Chiseled {
		return "minecraft:chiseled_resin_bricks", 0
	}
	return "minecraft:resin_bricks", 0
}

func (r ResinBricks) EncodeBlock() (string, map[string]any) {
	if r.Chiseled {
		return "minecraft:chiseled_resin_bricks", nil
	}
	return "minecraft:resin_bricks", nil
}
