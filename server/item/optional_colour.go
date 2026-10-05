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

type OptionalColour uint8

var colours = Colours()

func NewOptionalColour(c Colour) OptionalColour {
	return OptionalColour(c.Uint8() + 1)
}

func (oc OptionalColour) Colour() (Colour, bool) {
	if oc == 0 {
		return Colour{}, false
	}
	return colours[(oc - 1)], true
}

func (oc OptionalColour) Uint8() uint8 {
	return uint8(oc)
}

func (oc OptionalColour) Prepend(str string) string {
	if oc != 0 {
		return colours[(oc-1)].String() + "_" + str
	}
	return str
}

func OptionalColours() []OptionalColour {
	optionalColours := make([]OptionalColour, 17)
	for i, c := range colours {
		optionalColours[i+1] = NewOptionalColour(c)
	}
	return optionalColours
}
