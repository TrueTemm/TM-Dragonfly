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

type SandstoneType struct {
	sandstone
}

type sandstone uint8

func NormalSandstone() SandstoneType {
	return SandstoneType{0}
}

func CutSandstone() SandstoneType {
	return SandstoneType{1}
}

func ChiseledSandstone() SandstoneType {
	return SandstoneType{2}
}

func SmoothSandstone() SandstoneType {
	return SandstoneType{3}
}

func (s sandstone) Uint8() uint8 {
	return uint8(s)
}

func (s sandstone) Name() string {
	switch s {
	case 0:
		return "Sandstone"
	case 1:
		return "Cut Sandstone"
	case 2:
		return "Chiseled Sandstone"
	case 3:
		return "Smooth Sandstone"
	}
	panic("unknown sandstone type")
}

func (s sandstone) String() string {
	switch s {
	case 0:
		return "default"
	case 1:
		return "cut"
	case 2:
		return "chiseled"
	case 3:
		return "smooth"
	}
	panic("unknown sandstone type")
}

func SandstoneTypes() []SandstoneType {
	return []SandstoneType{NormalSandstone(), CutSandstone(), ChiseledSandstone(), SmoothSandstone()}
}
