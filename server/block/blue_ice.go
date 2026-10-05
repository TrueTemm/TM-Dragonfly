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

type BlueIce struct {
	solid
}

func (b BlueIce) BreakInfo() BreakInfo {
	return newBreakInfo(2.8, alwaysHarvestable, pickaxeEffective, silkTouchOnlyDrop(b))
}

func (b BlueIce) Friction() float64 {
	return 0.989
}

func (BlueIce) EncodeItem() (name string, meta int16) {
	return "minecraft:blue_ice", 0
}

func (BlueIce) EncodeBlock() (string, map[string]any) {
	return "minecraft:blue_ice", nil
}
