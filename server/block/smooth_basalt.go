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

type SmoothBasalt struct {
	solid
	bassDrum
}

func (SmoothBasalt) EncodeBlock() (string, map[string]any) {
	return "minecraft:smooth_basalt", nil
}

func (SmoothBasalt) EncodeItem() (name string, meta int16) {
	return "minecraft:smooth_basalt", 0
}

func (s SmoothBasalt) BreakInfo() BreakInfo {
	return newBreakInfo(1.25, pickaxeHarvestable, pickaxeEffective, oneOf(s)).withBlastResistance(4.2)
}
