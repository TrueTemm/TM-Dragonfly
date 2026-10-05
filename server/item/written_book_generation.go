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

type WrittenBookGeneration struct {
	generation
}

type generation uint8

func OriginalGeneration() WrittenBookGeneration {
	return WrittenBookGeneration{0}
}

func CopyGeneration() WrittenBookGeneration {
	return WrittenBookGeneration{1}
}

func CopyOfCopyGeneration() WrittenBookGeneration {
	return WrittenBookGeneration{2}
}

func (g generation) Uint8() uint8 {
	return uint8(g)
}

func (g generation) String() string {
	switch g {
	case 0:
		return "original"
	case 1:
		return "copy of original"
	case 2:
		return "copy of copy"
	}
	panic("unknown written book generation")
}
