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

type SkullType struct {
	skull
}

func SkeletonSkull() SkullType {
	return SkullType{0}
}

func WitherSkeletonSkull() SkullType {
	return SkullType{1}
}

func ZombieHead() SkullType {
	return SkullType{2}
}

func PlayerHead() SkullType {
	return SkullType{3}
}

func CreeperHead() SkullType {
	return SkullType{4}
}

func DragonHead() SkullType {
	return SkullType{5}
}

func PiglinHead() SkullType {
	return SkullType{6}
}

func SkullTypes() []SkullType {
	return []SkullType{SkeletonSkull(), WitherSkeletonSkull(), ZombieHead(), PlayerHead(), CreeperHead(), DragonHead(), PiglinHead()}
}

type skull uint8

func (s skull) Uint8() uint8 {
	return uint8(s)
}

func (s skull) Name() string {
	switch s {
	case 0:
		return "Skeleton Skull"
	case 1:
		return "Wither Skeleton Skull"
	case 2:
		return "Zombie Head"
	case 3:
		return "Player Head"
	case 4:
		return "Creeper Head"
	case 5:
		return "Dragon Head"
	case 6:
		return "Piglin Head"
	}
	panic("unknown skull type")
}

func (s skull) String() string {
	switch s {
	case 0:
		return "skeleton_skull"
	case 1:
		return "wither_skeleton_skull"
	case 2:
		return "zombie_head"
	case 3:
		return "player_head"
	case 4:
		return "creeper_head"
	case 5:
		return "dragon_head"
	case 6:
		return "piglin_head"
	}
	panic("unknown skull type")
}
