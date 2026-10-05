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

type Shroomlight struct {
	solid
}

func (Shroomlight) LightEmissionLevel() uint8 {
	return 15
}

func (s Shroomlight) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, hoeEffective, oneOf(s))
}

func (Shroomlight) CompostChance() float64 {
	return 0.65
}

func (Shroomlight) EncodeItem() (name string, meta int16) {
	return "minecraft:shroomlight", 0
}

func (Shroomlight) EncodeBlock() (string, map[string]any) {
	return "minecraft:shroomlight", nil
}
