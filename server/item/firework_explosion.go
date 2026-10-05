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

type FireworkExplosion struct {
	Shape FireworkShape

	Colour Colour

	Fade Colour

	Fades bool

	Twinkle bool

	Trail bool
}

func (f FireworkExplosion) EncodeNBT() map[string]any {
	data := map[string]any{
		"FireworkType":    f.Shape.Uint8(),
		"FireworkColor":   [1]uint8{uint8(invertColour(f.Colour))},
		"FireworkFade":    [0]uint8{},
		"FireworkFlicker": boolByte(f.Twinkle),
		"FireworkTrail":   boolByte(f.Trail),
	}
	if f.Fades {
		data["FireworkFade"] = [1]uint8{uint8(invertColour(f.Fade))}
	}
	return data
}

func (f FireworkExplosion) DecodeNBT(data map[string]any) any {
	f.Shape = FireworkShapes()[data["FireworkType"].(uint8)]
	f.Twinkle = data["FireworkFlicker"].(uint8) == 1
	f.Trail = data["FireworkTrail"].(uint8) == 1

	colours := data["FireworkColor"]
	if diskColour, ok := colours.([1]uint8); ok {
		f.Colour = invertColourID(int16(diskColour[0]))
	} else if networkColours, ok := colours.([]any); ok {
		f.Colour = invertColourID(int16(networkColours[0].(uint8)))
	}

	if fades, ok := data["FireworkFade"].([1]uint8); ok {
		f.Fade, f.Fades = invertColourID(int16(fades[0])), true
	}
	return f
}
