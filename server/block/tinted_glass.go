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

type TintedGlass struct {
	solid
	clicksAndSticks
}

func (TintedGlass) PreventsSuffocation() bool {
	return true
}

func (g TintedGlass) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, oneOf(g)).withBlastResistance(0.3)
}

func (TintedGlass) EncodeItem() (name string, meta int16) {
	return "minecraft:tinted_glass", 0
}

func (TintedGlass) EncodeBlock() (string, map[string]any) {
	return "minecraft:tinted_glass", nil
}
