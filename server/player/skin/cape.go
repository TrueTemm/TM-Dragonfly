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

package skin

import (
	"image"
	"image/color"
)

type Cape struct {
	w, h int

	Pix []uint8
}

func NewCape(width, height int) Cape {
	return Cape{w: width, h: height, Pix: make([]uint8, width*height*4)}
}

func (c Cape) ColorModel() color.Model {
	return color.RGBAModel
}

func (c Cape) Bounds() image.Rectangle {
	return image.Rectangle{
		Max: image.Point{X: c.w, Y: c.h},
	}
}

func (c Cape) At(x, y int) color.Color {
	if x < 0 || y < 0 || x >= c.w || y >= c.h {
		panic("pixel coordinates out of bounds")
	}
	offset := x*4 + c.w*y*4
	return color.RGBA{
		R: c.Pix[offset],
		G: c.Pix[offset+1],
		B: c.Pix[offset+2],
		A: c.Pix[offset+3],
	}
}
