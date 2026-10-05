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

package v729

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v748"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

func applyDeltas729(p packet.Pool) {
	delete(p, v786.IDMovementEffect)
	delete(p, v786.IDSetMovementAuthority)
	p[v786.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	p[v786.IDInventoryContent] = func() packet.Packet { return &InventoryContent{} }
	p[v786.IDInventorySlot] = func() packet.Packet { return &InventorySlot{} }
	p[v786.IDMobEffect] = func() packet.Packet { return &MobEffect{} }
	p[v786.IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }
}

func NewClientPool() packet.Pool {
	p := v748.NewClientPool()
	applyDeltas729(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v748.NewServerPool()
	applyDeltas729(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v748.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v748.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.StartGame:
		out, _ := v748.FromLatestShared(proto, pk)
		sg := out[0].(*v766.StartGame)
		sg.BaseGameVersion, sg.GameVersion = "1.21.30", "1.21.30"
		items := itemdata.Items729()
		sg.Items = make([]v766.ItemEntry, len(items))
		for i, it := range items {
			sg.Items[i] = v766.ItemEntry{Name: it.Name, RuntimeID: it.RuntimeID, ComponentBased: it.ComponentBased}
		}
		return []packet.Packet{sg}, true
	case *packet.PlayerAuthInput:
		return []packet.Packet{fromLatestPlayerAuthInput(pk)}, true
	case *packet.InventoryContent:
		return []packet.Packet{fromLatestInventoryContent(pk)}, true
	case *packet.InventorySlot:
		return []packet.Packet{fromLatestInventorySlot(pk)}, true
	case *packet.MobEffect:
		return []packet.Packet{fromLatestMobEffect(pk)}, true
	case *packet.ResourcePacksInfo:
		return []packet.Packet{fromLatestResourcePacksInfo(pk)}, true
	case *packet.MovementEffect:
		return nil, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *PlayerAuthInput:
		return []packet.Packet{toLatestPlayerAuthInput(pk)}, true
	case *InventoryContent:
		return []packet.Packet{toLatestInventoryContent(pk)}, true
	case *InventorySlot:
		return []packet.Packet{toLatestInventorySlot(pk)}, true
	case *MobEffect:
		return []packet.Packet{toLatestMobEffect(pk)}, true
	case *ResourcePacksInfo:
		return []packet.Packet{&packet.ResourcePacksInfo{}}, true
	}
	return nil, false
}

func FromLatestResourcePacksInfo(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	return fromLatestResourcePacksInfo(pk)
}

func ToLatestPlayerAuthInput(pk *PlayerAuthInput) *packet.PlayerAuthInput {
	return toLatestPlayerAuthInput(pk)
}
func FromLatestPlayerAuthInput(pk *packet.PlayerAuthInput) *PlayerAuthInput {
	return fromLatestPlayerAuthInput(pk)
}
