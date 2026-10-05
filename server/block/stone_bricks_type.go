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

type StoneBricksType struct {
	stoneBricks
}

type stoneBricks uint8

func NormalStoneBricks() StoneBricksType {
	return StoneBricksType{0}
}

func MossyStoneBricks() StoneBricksType {
	return StoneBricksType{1}
}

func CrackedStoneBricks() StoneBricksType {
	return StoneBricksType{2}
}

func ChiseledStoneBricks() StoneBricksType {
	return StoneBricksType{3}
}

func (s stoneBricks) Uint8() uint8 {
	return uint8(s)
}

func (s stoneBricks) Name() string {
	switch s {
	case 0:
		return "Stone Bricks"
	case 1:
		return "Mossy Stone Bricks"
	case 2:
		return "Cracked Stone Bricks"
	case 3:
		return "Chiseled Stone Bricks"
	}
	panic("unknown stone bricks type")
}

func (s stoneBricks) String() string {
	switch s {
	case 0:
		return "stone_bricks"
	case 1:
		return "mossy_stone_bricks"
	case 2:
		return "cracked_stone_bricks"
	case 3:
		return "chiseled_stone_bricks"
	}
	panic("unknown stone bricks type")
}

func StoneBricksTypes() []StoneBricksType {
	return []StoneBricksType{NormalStoneBricks(), MossyStoneBricks(), CrackedStoneBricks(), ChiseledStoneBricks()}
}
