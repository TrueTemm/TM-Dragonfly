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

type SmithingTemplateType struct {
	smithingTemplateType
}

func TemplateNetheriteUpgrade() SmithingTemplateType {
	return SmithingTemplateType{0}
}

func TemplateSentry() SmithingTemplateType {
	return SmithingTemplateType{1}
}

func TemplateVex() SmithingTemplateType {
	return SmithingTemplateType{2}
}

func TemplateWild() SmithingTemplateType {
	return SmithingTemplateType{3}
}

func TemplateCoast() SmithingTemplateType {
	return SmithingTemplateType{4}
}

func TemplateDune() SmithingTemplateType {
	return SmithingTemplateType{5}
}

func TemplateWayFinder() SmithingTemplateType {
	return SmithingTemplateType{6}
}

func TemplateRaiser() SmithingTemplateType {
	return SmithingTemplateType{7}
}

func TemplateShaper() SmithingTemplateType {
	return SmithingTemplateType{8}
}

func TemplateHost() SmithingTemplateType {
	return SmithingTemplateType{9}
}

func TemplateWard() SmithingTemplateType {
	return SmithingTemplateType{10}
}

func TemplateSilence() SmithingTemplateType {
	return SmithingTemplateType{11}
}

func TemplateTide() SmithingTemplateType {
	return SmithingTemplateType{12}
}

func TemplateSnout() SmithingTemplateType {
	return SmithingTemplateType{13}
}

func TemplateRib() SmithingTemplateType {
	return SmithingTemplateType{14}
}

func TemplateEye() SmithingTemplateType {
	return SmithingTemplateType{15}
}

func TemplateSpire() SmithingTemplateType {
	return SmithingTemplateType{16}
}

func TemplateFlow() SmithingTemplateType {
	return SmithingTemplateType{17}
}

func TemplateBolt() SmithingTemplateType {
	return SmithingTemplateType{18}
}

func SmithingTemplates() []SmithingTemplateType {
	return []SmithingTemplateType{
		TemplateNetheriteUpgrade(),
		TemplateSentry(),
		TemplateVex(),
		TemplateWild(),
		TemplateCoast(),
		TemplateDune(),
		TemplateWayFinder(),
		TemplateRaiser(),
		TemplateShaper(),
		TemplateHost(),
		TemplateWard(),
		TemplateSilence(),
		TemplateTide(),
		TemplateSnout(),
		TemplateRib(),
		TemplateEye(),
		TemplateSpire(),
		TemplateFlow(),
		TemplateBolt(),
	}
}

type smithingTemplateType uint8

func (s smithingTemplateType) String() string {
	switch s {
	case 0:
		return "netherite_upgrade"
	case 1:
		return "sentry"
	case 2:
		return "vex"
	case 3:
		return "wild"
	case 4:
		return "coast"
	case 5:
		return "dune"
	case 6:
		return "wayfinder"
	case 7:
		return "raiser"
	case 8:
		return "shaper"
	case 9:
		return "host"
	case 10:
		return "ward"
	case 11:
		return "silence"
	case 12:
		return "tide"
	case 13:
		return "snout"
	case 14:
		return "rib"
	case 15:
		return "eye"
	case 16:
		return "spire"
	case 17:
		return "flow"
	case 18:
		return "bolt"
	}

	panic("should never happen")
}

func smithingTemplateFromString(name string) (SmithingTemplateType, bool) {
	switch name {
	case "netherite_upgrade":
		return TemplateNetheriteUpgrade(), true
	case "sentry":
		return TemplateSentry(), true
	case "vex":
		return TemplateVex(), true
	case "wild":
		return TemplateWild(), true
	case "coast":
		return TemplateCoast(), true
	case "dune":
		return TemplateDune(), true
	case "wayfinder":
		return TemplateWayFinder(), true
	case "raiser":
		return TemplateRaiser(), true
	case "shaper":
		return TemplateShaper(), true
	case "host":
		return TemplateHost(), true
	case "ward":
		return TemplateWard(), true
	case "silence":
		return TemplateSilence(), true
	case "tide":
		return TemplateTide(), true
	case "snout":
		return TemplateSnout(), true
	case "rib":
		return TemplateRib(), true
	case "eye":
		return TemplateEye(), true
	case "spire":
		return TemplateSpire(), true
	case "flow":
		return TemplateFlow(), true
	case "bolt":
		return TemplateBolt(), true
	default:
		return SmithingTemplateType{}, false
	}
}
