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

type BlackstoneType struct {
	blackstone
}

type blackstone uint8

func NormalBlackstone() BlackstoneType {
	return BlackstoneType{0}
}

func GildedBlackstone() BlackstoneType {
	return BlackstoneType{1}
}

func PolishedBlackstone() BlackstoneType {
	return BlackstoneType{2}
}

func ChiseledPolishedBlackstone() BlackstoneType {
	return BlackstoneType{3}
}

func (s blackstone) Uint8() uint8 {
	return uint8(s)
}

func (s blackstone) Name() string {
	switch s {
	case 0:
		return "Blackstone"
	case 1:
		return "Gilded Blackstone"
	case 2:
		return "Polished Blackstone"
	case 3:
		return "Chiseled Polished Blackstone"
	}
	panic("unknown blackstone type")
}

func (s blackstone) String() string {
	switch s {
	case 0:
		return "blackstone"
	case 1:
		return "gilded_blackstone"
	case 2:
		return "polished_blackstone"
	case 3:
		return "chiseled_polished_blackstone"
	}
	panic("unknown blackstone type")
}

func BlackstoneTypes() []BlackstoneType {
	return []BlackstoneType{NormalBlackstone(), GildedBlackstone(), PolishedBlackstone(), ChiseledPolishedBlackstone()}
}
