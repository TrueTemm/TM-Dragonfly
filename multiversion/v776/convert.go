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

package v776

import (
	_ "embed"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

//go:embed biome_definitions_554.dat
var biomeDefinitions554 []byte

func BiomeDefinitions554() []byte { return biomeDefinitions554 }

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v786.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v786.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.LevelSoundEvent:
		return []packet.Packet{fromLatestLevelSoundEvent(proto, pk)}, true
	case *packet.ClientMovementPredictionSync:
		return []packet.Packet{fromLatestClientMovementPredictionSync(pk)}, true
	case *packet.SetHud:
		return []packet.Packet{fromLatestSetHud(pk)}, true
	case *packet.StartGame:

		sg := v786.FromLatestStartGame786(pk)
		sg.BaseGameVersion, sg.GameVersion = "1.21.60", "1.21.60"
		return []packet.Packet{sg}, true
	case *packet.ItemRegistry:

		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items776()}}, true
	case *packet.BiomeDefinitionList:
		return []packet.Packet{&v786.BiomeDefinitionList{SerialisedBiomeDefinitions: biomeDefinitions554}}, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *LevelSoundEvent:
		return []packet.Packet{toLatestLevelSoundEvent(proto, pk)}, true
	case *ClientMovementPredictionSync:
		return []packet.Packet{toLatestClientMovementPredictionSync(pk)}, true
	case *SetHud:
		return []packet.Packet{toLatestSetHud(pk)}, true
	}
	return nil, false
}
