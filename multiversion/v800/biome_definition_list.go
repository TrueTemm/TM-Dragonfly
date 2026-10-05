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

package v800

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type BiomeDefinitionList struct {
	BiomeDefinitions []BiomeDefinition
	StringList       []string
}

func (*BiomeDefinitionList) ID() uint32 { return packet.IDBiomeDefinitionList }

func (pk *BiomeDefinitionList) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.BiomeDefinitions)
	protocol.FuncSlice(io, &pk.StringList, io.String)
}

type BiomeDefinition struct {
	NameIndex        int16
	BiomeID          protocol.Optional[uint16]
	Temperature      float32
	Downfall         float32
	RedSporeDensity  float32
	BlueSporeDensity float32
	AshDensity       float32
	WhiteAshDensity  float32
	Depth            float32
	Scale            float32
	MapWaterColour   int32
	Rain             bool
	Tags             protocol.Optional[[]uint16]
}

func (x *BiomeDefinition) Marshal(io protocol.IO) {
	io.Int16(&x.NameIndex)
	protocol.OptionalFunc(io, &x.BiomeID, io.Uint16)
	io.Float32(&x.Temperature)
	io.Float32(&x.Downfall)
	io.Float32(&x.RedSporeDensity)
	io.Float32(&x.BlueSporeDensity)
	io.Float32(&x.AshDensity)
	io.Float32(&x.WhiteAshDensity)
	io.Float32(&x.Depth)
	io.Float32(&x.Scale)
	io.Int32(&x.MapWaterColour)
	io.Bool(&x.Rain)
	protocol.OptionalFunc(io, &x.Tags, func(s *[]uint16) {
		protocol.FuncSlice(io, s, io.Uint16)
	})

	var chunkGenerationPresent bool
	io.Bool(&chunkGenerationPresent)
}

func fromLatestBiomeDefinitionList800(pk *packet.BiomeDefinitionList) *BiomeDefinitionList {
	defs := make([]BiomeDefinition, len(pk.BiomeDefinitions))
	for i, d := range pk.BiomeDefinitions {

		var id protocol.Optional[uint16]
		if d.BiomeID >= 0 {
			id = protocol.Option(uint16(d.BiomeID))
		}
		defs[i] = BiomeDefinition{
			NameIndex:      d.NameIndex,
			BiomeID:        id,
			Temperature:    d.Temperature,
			Downfall:       d.Downfall,
			Depth:          d.Depth,
			Scale:          d.Scale,
			MapWaterColour: d.MapWaterColour,
			Rain:           d.Rain,
			Tags:           d.Tags,
		}
	}
	return &BiomeDefinitionList{BiomeDefinitions: defs, StringList: pk.StringList}
}
