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

type InfestedCobblestone struct {
	solid
	flute
}

func (i InfestedCobblestone) BreakInfo() BreakInfo {
	return newBreakInfo(1, pickaxeHarvestable, pickaxeEffective, silkTouchOnlyDrop(i)).withBlastResistance(0.75)
}

func (InfestedCobblestone) EncodeItem() (name string, meta int16) {
	return "minecraft:infested_cobblestone", 0
}

func (InfestedCobblestone) EncodeBlock() (string, map[string]any) {
	return "minecraft:infested_cobblestone", nil
}
