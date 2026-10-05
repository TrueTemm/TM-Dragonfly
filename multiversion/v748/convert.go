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

package v748

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

func applyDeltas748(p packet.Pool) {
	delete(p, v786.IDCameraAimAssistPresets)
	p[v786.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	p[v786.IDPlayerArmourDamage] = func() packet.Packet { return &PlayerArmourDamage{} }
	p[v786.IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }
}

func NewClientPool() packet.Pool {
	p := v766.NewClientPool()
	applyDeltas748(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v766.NewServerPool()
	applyDeltas748(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v766.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v766.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.StartGame:

		out, _ := v766.FromLatestShared(proto, pk)
		sg := out[0].(*v766.StartGame)
		sg.BaseGameVersion, sg.GameVersion = "1.21.40", "1.21.40"
		items := itemdata.Items748()
		sg.Items = make([]v766.ItemEntry, len(items))
		for i, it := range items {
			sg.Items[i] = v766.ItemEntry{Name: it.Name, RuntimeID: it.RuntimeID, ComponentBased: it.ComponentBased}
		}
		return []packet.Packet{sg}, true
	case *packet.PlayerAuthInput:
		return []packet.Packet{fromLatestPlayerAuthInput(pk)}, true
	case *packet.PlayerArmourDamage:
		return []packet.Packet{fromLatestPlayerArmourDamage(pk)}, true
	case *packet.ResourcePacksInfo:
		return []packet.Packet{fromLatestResourcePacksInfo(pk)}, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *PlayerAuthInput:
		return []packet.Packet{toLatestPlayerAuthInput(pk)}, true
	case *PlayerArmourDamage:
		return []packet.Packet{toLatestPlayerArmourDamage(pk)}, true
	case *ResourcePacksInfo:
		return []packet.Packet{&packet.ResourcePacksInfo{}}, true
	}
	return nil, false
}

func ToLatestPlayerAuthInput(pk *PlayerAuthInput) *packet.PlayerAuthInput {
	return toLatestPlayerAuthInput(pk)
}
func FromLatestPlayerAuthInput(pk *packet.PlayerAuthInput) *PlayerAuthInput {
	return fromLatestPlayerAuthInput(pk)
}
func FromLatestResourcePacksInfo(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	return fromLatestResourcePacksInfo(pk)
}
