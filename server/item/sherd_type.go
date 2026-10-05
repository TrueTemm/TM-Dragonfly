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

type SherdType struct {
	sherdType
}

func SherdTypeAngler() SherdType {
	return SherdType{0}
}

func SherdTypeArcher() SherdType {
	return SherdType{1}
}

func SherdTypeArmsUp() SherdType {
	return SherdType{2}
}

func SherdTypeBlade() SherdType {
	return SherdType{3}
}

func SherdTypeBrewer() SherdType {
	return SherdType{4}
}

func SherdTypeBurn() SherdType {
	return SherdType{5}
}

func SherdTypeDanger() SherdType {
	return SherdType{6}
}

func SherdTypeExplorer() SherdType {
	return SherdType{7}
}

func SherdTypeFriend() SherdType {
	return SherdType{8}
}

func SherdTypeHeart() SherdType {
	return SherdType{9}
}

func SherdTypeHeartbreak() SherdType {
	return SherdType{10}
}

func SherdTypeHowl() SherdType {
	return SherdType{11}
}

func SherdTypeMiner() SherdType {
	return SherdType{12}
}

func SherdTypeMourner() SherdType {
	return SherdType{13}
}

func SherdTypePlenty() SherdType {
	return SherdType{14}
}

func SherdTypePrize() SherdType {
	return SherdType{15}
}

func SherdTypeSheaf() SherdType {
	return SherdType{16}
}

func SherdTypeShelter() SherdType {
	return SherdType{17}
}

func SherdTypeSkull() SherdType {
	return SherdType{18}
}

func SherdTypeSnort() SherdType {
	return SherdType{19}
}

func SherdTypeFlow() SherdType {
	return SherdType{20}
}

func SherdTypeGuster() SherdType {
	return SherdType{21}
}

func SherdTypeScrape() SherdType {
	return SherdType{22}
}

func SherdTypes() []SherdType {
	return []SherdType{
		SherdTypeAngler(), SherdTypeArcher(), SherdTypeArmsUp(), SherdTypeBlade(), SherdTypeBrewer(), SherdTypeBurn(),
		SherdTypeDanger(), SherdTypeExplorer(), SherdTypeFriend(), SherdTypeHeart(), SherdTypeHeartbreak(), SherdTypeHowl(),
		SherdTypeMiner(), SherdTypeMourner(), SherdTypePlenty(), SherdTypePrize(), SherdTypeSheaf(), SherdTypeShelter(),
		SherdTypeSkull(), SherdTypeSnort(), SherdTypeFlow(), SherdTypeGuster(), SherdTypeScrape(),
	}
}

type sherdType uint8

func (c sherdType) String() string {
	switch c {
	case 0:
		return "angler"
	case 1:
		return "archer"
	case 2:
		return "arms_up"
	case 3:
		return "blade"
	case 4:
		return "brewer"
	case 5:
		return "burn"
	case 6:
		return "danger"
	case 7:
		return "explorer"
	case 8:
		return "friend"
	case 9:
		return "heart"
	case 10:
		return "heartbreak"
	case 11:
		return "howl"
	case 12:
		return "miner"
	case 13:
		return "mourner"
	case 14:
		return "plenty"
	case 15:
		return "prize"
	case 16:
		return "sheaf"
	case 17:
		return "shelter"
	case 18:
		return "skull"
	case 19:
		return "snort"
	case 20:
		return "flow"
	case 21:
		return "guster"
	case 22:
		return "scrape"
	}
	panic("unknown sherd type")
}

func (c sherdType) Uint8() uint8 {
	return uint8(c)
}
