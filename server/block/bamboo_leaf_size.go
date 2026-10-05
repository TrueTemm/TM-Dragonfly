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

type BambooLeafSize struct {
	bamboo
}

type bamboo uint8

func BambooSizeNoLeaves() BambooLeafSize {
	return BambooLeafSize{0}
}

func BambooSizeSmallLeaves() BambooLeafSize {
	return BambooLeafSize{1}
}

func BambooSizeLargeLeaves() BambooLeafSize {
	return BambooLeafSize{2}
}

func (b bamboo) Uint8() uint8 {
	return uint8(b)
}

func (b bamboo) String() string {
	switch b {
	case 0:
		return "no_leaves"
	case 1:
		return "small_leaves"
	case 2:
		return "large_leaves"
	}
	panic("unknown bamboo leaf size")
}

func (b bamboo) Name() string {
	switch b {
	case 0:
		return "No Leaves"
	case 1:
		return "Small Leaves"
	case 2:
		return "Large Leaves"
	}
	panic("unknown bamboo leaf size")
}

func BambooLeafSizes() []BambooLeafSize {
	return []BambooLeafSize{BambooSizeNoLeaves(), BambooSizeSmallLeaves(), BambooSizeLargeLeaves()}
}
