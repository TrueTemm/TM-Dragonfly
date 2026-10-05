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

type AnvilType struct {
	anvil
}

func UndamagedAnvil() AnvilType {
	return AnvilType{0}
}

func SlightlyDamagedAnvil() AnvilType {
	return AnvilType{1}
}

func VeryDamagedAnvil() AnvilType {
	return AnvilType{2}
}

func AnvilTypes() []AnvilType {
	return []AnvilType{UndamagedAnvil(), SlightlyDamagedAnvil(), VeryDamagedAnvil()}
}

type anvil uint8

func (a anvil) Uint8() uint8 {
	return uint8(a)
}

func (a anvil) String() string {
	switch a {
	case 0:
		return "anvil"
	case 1:
		return "chipped_anvil"
	case 2:
		return "damaged_anvil"
	}
	panic("should never happen")
}
