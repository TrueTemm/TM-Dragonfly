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

import (
	"image/color"
)

type Colour struct {
	colour
}

func ColourWhite() Colour {
	return Colour{0}
}

func ColourOrange() Colour {
	return Colour{1}
}

func ColourMagenta() Colour {
	return Colour{2}
}

func ColourLightBlue() Colour {
	return Colour{3}
}

func ColourYellow() Colour {
	return Colour{4}
}

func ColourLime() Colour {
	return Colour{5}
}

func ColourPink() Colour {
	return Colour{6}
}

func ColourGrey() Colour {
	return Colour{7}
}

func ColourLightGrey() Colour {
	return Colour{8}
}

func ColourCyan() Colour {
	return Colour{9}
}

func ColourPurple() Colour {
	return Colour{10}
}

func ColourBlue() Colour {
	return Colour{11}
}

func ColourBrown() Colour {
	return Colour{12}
}

func ColourGreen() Colour {
	return Colour{13}
}

func ColourRed() Colour {
	return Colour{14}
}

func ColourBlack() Colour {
	return Colour{15}
}

func Colours() []Colour {
	return []Colour{
		ColourWhite(), ColourOrange(), ColourMagenta(), ColourLightBlue(), ColourYellow(), ColourLime(), ColourPink(), ColourGrey(),
		ColourLightGrey(), ColourCyan(), ColourPurple(), ColourBlue(), ColourBrown(), ColourGreen(), ColourRed(), ColourBlack(),
	}
}

type colour uint8

func (c colour) RGBA() color.RGBA {
	switch c {
	case 0:
		return color.RGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff}
	case 1:
		return color.RGBA{R: 0xf9, G: 0x80, B: 0x1d, A: 0xff}
	case 2:
		return color.RGBA{R: 0xc7, G: 0x4e, B: 0xbd, A: 0xff}
	case 3:
		return color.RGBA{R: 0x3a, G: 0xb3, B: 0xda, A: 0xff}
	case 4:
		return color.RGBA{R: 0xfe, G: 0xd8, B: 0x3d, A: 0xff}
	case 5:
		return color.RGBA{R: 0x80, G: 0xc7, B: 0x1f, A: 0xff}
	case 6:
		return color.RGBA{R: 0xf3, G: 0x8b, B: 0xaa, A: 0xff}
	case 7:
		return color.RGBA{R: 0x47, G: 0x4f, B: 0x52, A: 0xff}
	case 8:
		return color.RGBA{R: 0x9d, G: 0x9d, B: 0x97, A: 0xff}
	case 9:
		return color.RGBA{R: 0x16, G: 0x9c, B: 0x9c, A: 0xff}
	case 10:
		return color.RGBA{R: 0x89, G: 0x32, B: 0xb8, A: 0xff}
	case 11:
		return color.RGBA{R: 0x3c, G: 0x44, B: 0xaa, A: 0xff}
	case 12:
		return color.RGBA{R: 0x83, G: 0x54, B: 0x32, A: 0xff}
	case 13:
		return color.RGBA{R: 0x5e, G: 0x7c, B: 0x16, A: 0xff}
	case 14:
		return color.RGBA{R: 0xb0, G: 0x2e, B: 0x26, A: 0xff}
	default:
		return color.RGBA{R: 0x1d, G: 0x1d, B: 0x21, A: 0xff}
	}
}

func (c colour) SignRGBA() color.RGBA {
	if c == 15 {
		return color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	}
	return c.RGBA()
}

func (c colour) String() string {
	switch c {
	default:
		return "white"
	case 1:
		return "orange"
	case 2:
		return "magenta"
	case 3:
		return "light_blue"
	case 4:
		return "yellow"
	case 5:
		return "lime"
	case 6:
		return "pink"
	case 7:
		return "gray"
	case 8:
		return "light_gray"
	case 9:
		return "cyan"
	case 10:
		return "purple"
	case 11:
		return "blue"
	case 12:
		return "brown"
	case 13:
		return "green"
	case 14:
		return "red"
	case 15:
		return "black"
	}
}

func (c colour) SilverString() string {
	if c == 8 {
		return "silver"
	}
	return c.String()
}

func (c colour) Uint8() uint8 {
	return uint8(c)
}

func invertColour(c Colour) int16 {
	return ^int16(c.Uint8()) & 0xf
}

func invertColourID(id int16) Colour {
	return Colours()[uint8(^id&0xf)]
}
