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

type GrindstoneAttachment struct {
	grindstoneAttachment
}

func StandingGrindstoneAttachment() GrindstoneAttachment {
	return GrindstoneAttachment{0}
}

func HangingGrindstoneAttachment() GrindstoneAttachment {
	return GrindstoneAttachment{1}
}

func WallGrindstoneAttachment() GrindstoneAttachment {
	return GrindstoneAttachment{2}
}

func GrindstoneAttachments() []GrindstoneAttachment {
	return []GrindstoneAttachment{StandingGrindstoneAttachment(), HangingGrindstoneAttachment(), WallGrindstoneAttachment()}
}

type grindstoneAttachment uint8

func (g grindstoneAttachment) Uint8() uint8 {
	return uint8(g)
}

func (g grindstoneAttachment) String() string {
	switch g {
	case 0:
		return "standing"
	case 1:
		return "hanging"
	case 2:
		return "side"
	}
	panic("should never happen")
}
