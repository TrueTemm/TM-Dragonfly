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

type Skin struct {
	w, h int

	Persona   bool
	PlayFabID string
	FullID    string

	Pix []uint8

	ModelConfig ModelConfig

	Model []byte

	Cape Cape

	Animations []Animation
}

func New(width, height int) Skin {
	return Skin{
		w:   width,
		h:   height,
		Pix: make([]uint8, width*height*4),
	}
}

func (s Skin) Bounds() image.Rectangle {
	return image.Rectangle{
		Max: image.Point{X: s.w, Y: s.h},
	}
}

func (s Skin) ColorModel() color.Model {
	return color.RGBAModel
}

func (s Skin) At(x, y int) color.Color {
	if x < 0 || y < 0 || x >= s.w || y >= s.h {
		panic("pixel coordinates out of bounds")
	}
	offset := x*4 + s.w*y*4
	return color.RGBA{
		R: s.Pix[offset],
		G: s.Pix[offset+1],
		B: s.Pix[offset+2],
		A: s.Pix[offset+3],
	}
}
