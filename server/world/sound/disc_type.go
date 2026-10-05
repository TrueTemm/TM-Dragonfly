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

package sound

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type DiscType struct {
	disc
}

func Disc13() DiscType {
	return DiscType{0}
}

func DiscCat() DiscType {
	return DiscType{1}
}

func DiscBlocks() DiscType {
	return DiscType{2}
}

func DiscChirp() DiscType {
	return DiscType{3}
}

func DiscFar() DiscType {
	return DiscType{4}
}

func DiscMall() DiscType {
	return DiscType{5}
}

func DiscMellohi() DiscType {
	return DiscType{6}
}

func DiscStal() DiscType {
	return DiscType{7}
}

func DiscStrad() DiscType {
	return DiscType{8}
}

func DiscWard() DiscType {
	return DiscType{9}
}

func Disc11() DiscType {
	return DiscType{10}
}

func DiscWait() DiscType {
	return DiscType{11}
}

func DiscOtherside() DiscType {
	return DiscType{12}
}

func DiscPigstep() DiscType {
	return DiscType{13}
}

func Disc5() DiscType {
	return DiscType{14}
}

func DiscRelic() DiscType {
	return DiscType{15}
}

func DiscCreator() DiscType {
	return DiscType{16}
}

func DiscCreatorMusicBox() DiscType {
	return DiscType{17}
}

func DiscPrecipice() DiscType {
	return DiscType{18}
}

func DiscTears() DiscType {
	return DiscType{19}
}

func DiscLavaChicken() DiscType {
	return DiscType{20}
}

func MusicDiscs() []DiscType {
	return []DiscType{
		Disc13(), DiscCat(), DiscBlocks(), DiscChirp(), DiscFar(), DiscMall(), DiscMellohi(), DiscStal(),
		DiscStrad(), DiscWard(), Disc11(), DiscWait(), DiscOtherside(), DiscPigstep(), Disc5(), DiscRelic(),
		DiscCreator(), DiscCreatorMusicBox(), DiscPrecipice(), DiscTears(), DiscLavaChicken(),
	}
}

type disc uint8

func (d disc) Uint8() uint8 {
	return uint8(d)
}

func (d disc) String() string {
	switch d {
	case 0:
		return "13"
	case 1:
		return "cat"
	case 2:
		return "blocks"
	case 3:
		return "chirp"
	case 4:
		return "far"
	case 5:
		return "mall"
	case 6:
		return "mellohi"
	case 7:
		return "stal"
	case 8:
		return "strad"
	case 9:
		return "ward"
	case 10:
		return "11"
	case 11:
		return "wait"
	case 12:
		return "otherside"
	case 13:
		return "pigstep"
	case 14:
		return "5"
	case 15:
		return "relic"
	case 16:
		return "creator"
	case 17:
		return "creator_music_box"
	case 18:
		return "precipice"
	case 19:
		return "tears"
	case 20:
		return "lava_chicken"
	}
	panic("unknown record type")
}

func (d disc) DisplayName() string {
	switch d {
	case 17:
		return "Creator (Music Box)"
	case 20:
		return "Lava Chicken"
	}
	if d > 12 {
		return cases.Title(language.English, cases.Compact).String(d.String())
	}
	return d.String()
}

func (d disc) Author() string {
	if d <= 11 {
		return "C418"
	}
	switch d {
	case 12, 13, 16, 17:
		return "Lena Raine"
	case 14:
		return "Samuel Åberg"
	case 15, 18:
		return "Aaron Cherof"
	case 19:
		return "Amos Roddy"
	case 20:
		return "Hyper Potions"
	}
	panic("unknown record type")
}
