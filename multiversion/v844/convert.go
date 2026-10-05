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

package v844

import (
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v827"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	return convertFromLatest(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	return convertToLatest(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	switch pk := pk.(type) {
	case *packet.BiomeDefinitionList:

		return []packet.Packet{pk}, true
	case *packet.Disconnect:

		return []packet.Packet{pk}, true
	case *packet.StartGame:
		return []packet.Packet{FromLatestStartGame844(pk)}, true
	case *packet.ItemRegistry:
		_ = pk
		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items844()}}, true
	case *packet.VoxelShapes:

		return nil, true
	}
	if out, ok := v827.FromLatestShared(proto, pk); ok {
		return out, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) (out []packet.Packet, ok bool) {
	if pk, isStartGame := pk.(*StartGame); isStartGame {
		return []packet.Packet{ToLatestStartGame844(pk)}, true
	}
	return v827.ToLatestShared(proto, pk)
}
