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

type TuffBricks struct {
	solid
	bassDrum

	Chiseled bool
}

func (t TuffBricks) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, pickaxeHarvestable, pickaxeEffective, oneOf(t)).withBlastResistance(6)
}

func (t TuffBricks) EncodeItem() (name string, meta int16) {
	if t.Chiseled {
		return "minecraft:chiseled_tuff_bricks", 0
	}
	return "minecraft:tuff_bricks", 0
}

func (t TuffBricks) EncodeBlock() (string, map[string]any) {
	if t.Chiseled {
		return "minecraft:chiseled_tuff_bricks", nil
	}
	return "minecraft:tuff_bricks", nil
}
