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

type OxidationType struct {
	oxidation
}

type oxidation uint8

func UnoxidisedOxidation() OxidationType {
	return OxidationType{0}
}

func ExposedOxidation() OxidationType {
	return OxidationType{1}
}

func WeatheredOxidation() OxidationType {
	return OxidationType{2}
}

func OxidisedOxidation() OxidationType {
	return OxidationType{3}
}

func (s oxidation) Uint8() uint8 {
	return uint8(s)
}

func (s oxidation) Name() string {
	switch s {
	case 0:
		return ""
	case 1:
		return "Exposed"
	case 2:
		return "Weathered"
	case 3:
		return "Oxidized"
	}
	panic("unknown oxidation type")
}

func (s oxidation) Decrease() (OxidationType, bool) {
	if s > 0 {
		return OxidationType{s - 1}, true
	}
	return UnoxidisedOxidation(), false
}

func (s oxidation) Increase() (OxidationType, bool) {
	if s < 3 {
		return OxidationType{s + 1}, true
	}
	return OxidisedOxidation(), false
}

func (s oxidation) String() string {
	switch s {
	case 0:
		return ""
	case 1:
		return "exposed"
	case 2:
		return "weathered"
	case 3:
		return "oxidized"
	}
	panic("unknown oxidation type")
}

func OxidationTypes() []OxidationType {
	return []OxidationType{UnoxidisedOxidation(), ExposedOxidation(), WeatheredOxidation(), OxidisedOxidation()}
}
