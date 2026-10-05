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

type CopperType struct {
	copper
}

type copper uint8

func NormalCopper() CopperType {
	return CopperType{0}
}

func CutCopper() CopperType {
	return CopperType{1}
}

func ChiseledCopper() CopperType {
	return CopperType{2}
}

func (s copper) Uint8() uint8 {
	return uint8(s)
}

func (s copper) Name() string {
	switch s {
	case 0:
		return "Copper"
	case 1:
		return "Cut Copper"
	case 2:
		return "Chiseled Copper"
	}
	panic("unknown copper type")
}

func (s copper) String() string {
	switch s {
	case 0:
		return "default"
	case 1:
		return "cut"
	case 2:
		return "chiseled"
	}
	panic("unknown copper type")
}

func CopperTypes() []CopperType {
	return []CopperType{NormalCopper(), CutCopper(), ChiseledCopper()}
}
