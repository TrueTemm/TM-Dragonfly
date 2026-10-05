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

type DeepslateType struct {
	deepslate
}

type deepslate uint8

func NormalDeepslate() DeepslateType {
	return DeepslateType{0}
}

func CobbledDeepslate() DeepslateType {
	return DeepslateType{1}
}

func PolishedDeepslate() DeepslateType {
	return DeepslateType{2}
}

func ChiseledDeepslate() DeepslateType {
	return DeepslateType{3}
}

func (s deepslate) Uint8() uint8 {
	return uint8(s)
}

func (s deepslate) Name() string {
	switch s {
	case 0:
		return "Deepslate"
	case 1:
		return "Cobbled Deepslate"
	case 2:
		return "Polished Deepslate"
	case 3:
		return "Chiseled Deepslate"
	}
	panic("unknown deepslate type")
}

func (s deepslate) String() string {
	switch s {
	case 0:
		return "deepslate"
	case 1:
		return "cobbled_deepslate"
	case 2:
		return "polished_deepslate"
	case 3:
		return "chiseled_deepslate"
	}
	panic("unknown deepslate type")
}

func DeepslateTypes() []DeepslateType {
	return []DeepslateType{NormalDeepslate(), CobbledDeepslate(), PolishedDeepslate(), ChiseledDeepslate()}
}
