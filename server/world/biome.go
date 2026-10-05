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

package world

import "image/color"

type Biome interface {
	Temperature() float64

	Rainfall() float64

	Depth() float64

	Scale() float64

	WaterColour() color.RGBA

	Tags() []string

	String() string

	EncodeBiome() int
}

var biomes = map[int]Biome{}

var biomeByName = map[string]Biome{}

func RegisterBiome(b Biome) {
	id := b.EncodeBiome()
	if _, ok := biomes[id]; ok {
		panic("cannot register the same biome (" + b.String() + ") twice")
	}
	biomes[id] = b
	biomeByName[b.String()] = b
}

func BiomeByID(id int) (Biome, bool) {
	e, ok := biomes[id]
	return e, ok
}

func BiomeByName(name string) (Biome, bool) {
	e, ok := biomeByName[name]
	return e, ok
}

func Biomes() []Biome {
	bs := make([]Biome, 0, len(biomes))
	for _, b := range biomes {
		bs = append(bs, b)
	}
	return bs
}

func ocean() Biome {
	if o, ok := BiomeByID(0); ok {
		return o
	}
	return unknownBiome{}
}

type unknownBiome struct {
	id int
}

func (unknownBiome) Temperature() float64    { return 0.5 }
func (unknownBiome) Rainfall() float64       { return 0 }
func (unknownBiome) Depth() float64          { return 0.1 }
func (unknownBiome) Scale() float64          { return 0.1 }
func (unknownBiome) WaterColour() color.RGBA { return color.RGBA{R: 0x44, G: 0xaf, B: 0xf5, A: 0xff} }
func (unknownBiome) Tags() []string          { return nil }
func (unknownBiome) String() string          { return "unknown" }
func (b unknownBiome) EncodeBiome() int      { return b.id }
