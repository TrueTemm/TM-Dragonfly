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

import (
	"github.com/df-mc/dragonfly/server/item"
)

type CoralType struct {
	coral
}

func TubeCoral() CoralType {
	return CoralType{0}
}

func BrainCoral() CoralType {
	return CoralType{1}
}

func BubbleCoral() CoralType {
	return CoralType{2}
}

func FireCoral() CoralType {
	return CoralType{3}
}

func HornCoral() CoralType {
	return CoralType{4}
}

func CoralTypes() []CoralType {
	return []CoralType{TubeCoral(), BrainCoral(), BubbleCoral(), FireCoral(), HornCoral()}
}

type coral uint8

func (c coral) Uint8() uint8 {
	return uint8(c)
}

func (c coral) Colour() item.Colour {
	switch c {
	case 0:
		return item.ColourBlue()
	case 1:
		return item.ColourPink()
	case 2:
		return item.ColourPurple()
	case 3:
		return item.ColourRed()
	case 4:
		return item.ColourYellow()
	}
	panic("unknown coral type")
}

func (c coral) Name() string {
	switch c {
	case 0:
		return "Tube Coral"
	case 1:
		return "Brain Coral"
	case 2:
		return "Bubble Coral"
	case 3:
		return "Fire Coral"
	case 4:
		return "Horn Coral"
	}
	panic("unknown coral type")
}

func (c coral) String() string {
	switch c {
	case 0:
		return "tube"
	case 1:
		return "brain"
	case 2:
		return "bubble"
	case 3:
		return "fire"
	case 4:
		return "horn"
	}
	panic("unknown coral type")
}
