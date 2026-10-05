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

type FireType struct {
	fire
}

type fire uint8

func NormalFire() FireType {
	return FireType{0}
}

func SoulFire() FireType {
	return FireType{1}
}

func (f fire) Uint8() uint8 {
	return uint8(f)
}

func (f fire) LightLevel() uint8 {
	switch f {
	case 0:
		return 15
	case 1:
		return 10
	}
	panic("unknown fire type")
}

func (f fire) Damage() float64 {
	switch f {
	case 0:
		return 1
	case 1:
		return 2
	}
	panic("unknown fire type")
}

func (f fire) Name() string {
	switch f {
	case 0:
		return "Fire"
	case 1:
		return "Soul Fire"
	}
	panic("unknown fire type")
}

func (f fire) String() string {
	switch f {
	case 0:
		return "normal"
	case 1:
		return "soul"
	}
	panic("unknown fire type")
}

func FireTypes() []FireType {
	return []FireType{NormalFire(), SoulFire()}
}
