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

type FireworkStar struct {
	FireworkExplosion
}

func (f FireworkStar) EncodeNBT() map[string]any {
	return map[string]any{
		"FireworksItem": f.FireworkExplosion.EncodeNBT(),
		"customColor":   int32FromRGBA(f.Colour.RGBA()),
	}
}

func (f FireworkStar) DecodeNBT(data map[string]any) any {
	if i, ok := data["FireworksItem"].(map[string]any); ok {
		f.FireworkExplosion = f.FireworkExplosion.DecodeNBT(i).(FireworkExplosion)
	}
	return f
}

func (f FireworkStar) EncodeItem() (name string, meta int16) {
	return "minecraft:firework_star", invertColour(f.Colour)
}
