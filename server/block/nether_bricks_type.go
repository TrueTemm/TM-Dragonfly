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

type NetherBricksType struct {
	netherBricks
}

type netherBricks uint8

func NormalNetherBricks() NetherBricksType {
	return NetherBricksType{0}
}

func RedNetherBricks() NetherBricksType {
	return NetherBricksType{1}
}

func CrackedNetherBricks() NetherBricksType {
	return NetherBricksType{2}
}

func ChiseledNetherBricks() NetherBricksType {
	return NetherBricksType{3}
}

func (n netherBricks) Uint8() uint8 {
	return uint8(n)
}

func (n netherBricks) Name() string {
	switch n {
	case 0:
		return "Nether Bricks"
	case 1:
		return "Red Nether Bricks"
	case 2:
		return "Cracked Nether Bricks"
	case 3:
		return "Chiseled Nether Bricks"
	}
	panic("unknown nether brick type")
}

func (n netherBricks) String() string {
	switch n {
	case 0:
		return "nether_brick"
	case 1:
		return "red_nether_brick"
	case 2:
		return "cracked_nether_bricks"
	case 3:
		return "chiseled_nether_bricks"
	}
	panic("unknown nether brick type")
}

func NetherBricksTypes() []NetherBricksType {
	return []NetherBricksType{NormalNetherBricks(), RedNetherBricks(), CrackedNetherBricks(), ChiseledNetherBricks()}
}
