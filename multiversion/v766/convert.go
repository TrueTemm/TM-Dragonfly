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

package v766

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v776"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

func applyDeltas766(p packet.Pool) {
	delete(p, v786.IDClientCameraAimAssist)
	delete(p, v786.IDClientMovementPredictionSync)
	p[v786.IDStartGame] = func() packet.Packet { return &StartGame{} }
	p[v786.IDItemRegistry] = func() packet.Packet { return &ItemComponent{} }
	p[v786.IDCreativeContent] = func() packet.Packet { return &CreativeContent{} }
	p[v786.IDUpdateAbilities] = func() packet.Packet { return &UpdateAbilities{} }
	p[v786.IDBossEvent] = func() packet.Packet { return &BossEvent{} }
	p[v786.IDCommandBlockUpdate] = func() packet.Packet { return &CommandBlockUpdate{} }
	p[v786.IDStructureBlockUpdate] = func() packet.Packet { return &StructureBlockUpdate{} }
}

func NewClientPool() packet.Pool {
	p := v776.NewClientPool()
	applyDeltas766(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v776.NewServerPool()
	applyDeltas766(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v776.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v776.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.StartGame:
		sg := v786.FromLatestStartGame786(pk)
		sg.BaseGameVersion, sg.GameVersion = "1.21.50", "1.21.50"
		items := itemdata.Items766()
		out := &StartGame{StartGame: *sg, Items: make([]ItemEntry, len(items))}
		for i, it := range items {
			out.Items[i] = ItemEntry{Name: it.Name, RuntimeID: it.RuntimeID, ComponentBased: it.ComponentBased}
		}
		return []packet.Packet{out}, true
	case *packet.ItemRegistry:

		return []packet.Packet{&ItemComponent{Items: []ItemComponentEntry{}}}, true
	case *packet.CreativeContent:
		return []packet.Packet{fromLatestCreativeContent(pk)}, true
	case *packet.UpdateAbilities:
		return []packet.Packet{fromLatestUpdateAbilities(pk)}, true
	case *packet.BossEvent:
		return []packet.Packet{fromLatestBossEvent(pk)}, true
	case *packet.CommandBlockUpdate:
		return []packet.Packet{fromLatestCommandBlockUpdate(pk)}, true
	case *packet.StructureBlockUpdate:
		return []packet.Packet{fromLatestStructureBlockUpdate(pk)}, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *UpdateAbilities:
		return []packet.Packet{toLatestUpdateAbilities(pk)}, true
	case *BossEvent:
		return []packet.Packet{toLatestBossEvent(pk)}, true
	case *CommandBlockUpdate:
		return []packet.Packet{toLatestCommandBlockUpdate(pk)}, true
	case *StructureBlockUpdate:
		return []packet.Packet{toLatestStructureBlockUpdate(pk)}, true
	case *CreativeContent:
		return []packet.Packet{toLatestCreativeContent(pk)}, true
	case *StartGame:
		return []packet.Packet{v786.ToLatestStartGame786(&pk.StartGame)}, true
	case *ItemComponent:
		return []packet.Packet{&packet.ItemRegistry{}}, true
	}
	return nil, false
}
