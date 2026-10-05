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

type ReinforcedDeepslate struct {
	solid
	bassDrum
}

func (r ReinforcedDeepslate) BreakInfo() BreakInfo {
	return newBreakInfo(55, alwaysHarvestable, nothingEffective, oneOf(r)).withBlastResistance(1200)
}

func (ReinforcedDeepslate) EncodeItem() (name string, meta int16) {
	return "minecraft:reinforced_deepslate", 0
}

func (ReinforcedDeepslate) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:reinforced_deepslate", nil
}
