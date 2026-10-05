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

package v686

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v712"
	"github.com/df-mc/dragonfly/multiversion/v729"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

func applyDeltas686(p packet.Pool) {
	delete(p, v786.IDServerBoundLoadingScreen)
	delete(p, v786.IDJigsawStructureData)
	delete(p, v786.IDCurrentStructureFeature)
	delete(p, v786.IDServerBoundDiagnostics)
	p[v786.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	p[v786.IDPlaySound] = func() packet.Packet { return &PlaySound{} }
	p[v786.IDChangeDimension] = func() packet.Packet { return &ChangeDimension{} }
	p[v786.IDDisconnect] = func() packet.Packet { return &Disconnect{} }
	p[v786.IDSetTitle] = func() packet.Packet { return &SetTitle{} }
	p[v786.IDStopSound] = func() packet.Packet { return &StopSound{} }
	p[v786.IDMobArmourEquipment] = func() packet.Packet { return &MobArmourEquipment{} }
	p[v786.IDPlayerArmourDamage] = func() packet.Packet { return &PlayerArmourDamage{} }
	p[v786.IDInventoryContent] = func() packet.Packet { return &InventoryContent{} }
	p[v786.IDInventorySlot] = func() packet.Packet { return &InventorySlot{} }
	p[v786.IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }
}

func NewClientPool() packet.Pool {
	p := v712.NewClientPool()
	applyDeltas686(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v712.NewServerPool()
	applyDeltas686(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v712.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v712.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.StartGame:
		out, _ := v712.FromLatestShared(proto, pk)
		sg := out[0].(*v766.StartGame)
		sg.BaseGameVersion, sg.GameVersion = "1.21.2", "1.21.2"
		items := itemdata.Items686()
		sg.Items = make([]v766.ItemEntry, len(items))
		for i, it := range items {
			sg.Items[i] = v766.ItemEntry{Name: it.Name, RuntimeID: it.RuntimeID, ComponentBased: it.ComponentBased}
		}
		return []packet.Packet{sg}, true
	case *packet.PlayerAuthInput:
		return []packet.Packet{&PlayerAuthInput{PlayerAuthInput: *v729.FromLatestPlayerAuthInput(pk)}}, true
	case *packet.PlaySound:
		return []packet.Packet{fromLatestPlaySound(pk)}, true
	case *packet.ChangeDimension:
		x := v786.FromLatestChangeDimension(pk)
		return []packet.Packet{&ChangeDimension{Dimension: x.Dimension, Position: x.Position, Respawn: x.Respawn}}, true
	case *packet.Disconnect:
		x := v786.FromLatestDisconnect(pk)
		return []packet.Packet{&Disconnect{Reason: x.Reason, HideDisconnectionScreen: x.HideDisconnectionScreen, Message: x.Message}}, true
	case *packet.SetTitle:
		return []packet.Packet{&SetTitle{ActionType: pk.ActionType, Text: pk.Text, FadeInDuration: pk.FadeInDuration,
			RemainDuration: pk.RemainDuration, FadeOutDuration: pk.FadeOutDuration, XUID: pk.XUID, PlatformOnlineID: pk.PlatformOnlineID}}, true
	case *packet.StopSound:
		return []packet.Packet{&StopSound{SoundName: pk.SoundName, StopAll: pk.StopAll}}, true
	case *packet.MobArmourEquipment:
		return []packet.Packet{&MobArmourEquipment{EntityRuntimeID: pk.EntityRuntimeID, Helmet: pk.Helmet,
			Chestplate: pk.Chestplate, Leggings: pk.Leggings, Boots: pk.Boots}}, true
	case *packet.PlayerArmourDamage:
		return []packet.Packet{fromLatestPlayerArmourDamage(pk)}, true
	case *packet.InventoryContent:
		return []packet.Packet{&InventoryContent{WindowID: pk.WindowID, Content: pk.Content}}, true
	case *packet.InventorySlot:
		return []packet.Packet{&InventorySlot{WindowID: pk.WindowID, Slot: pk.Slot, NewItem: pk.NewItem}}, true
	case *packet.ResourcePacksInfo:
		return []packet.Packet{fromLatestResourcePacksInfo(pk)}, true
	case *packet.JigsawStructureData, *packet.CurrentStructureFeature:
		return nil, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *PlayerAuthInput:
		return []packet.Packet{v729.ToLatestPlayerAuthInput(&pk.PlayerAuthInput)}, true
	case *PlaySound:
		return []packet.Packet{v786.ToLatestPlaySound(&v786.PlaySound{SoundName: pk.SoundName, Position: pk.Position, Volume: pk.Volume, Pitch: pk.Pitch})}, true
	case *SetTitle:
		return []packet.Packet{&packet.SetTitle{ActionType: pk.ActionType, Text: pk.Text, FadeInDuration: pk.FadeInDuration,
			RemainDuration: pk.RemainDuration, FadeOutDuration: pk.FadeOutDuration, XUID: pk.XUID, PlatformOnlineID: pk.PlatformOnlineID}}, true
	case *StopSound:
		return []packet.Packet{&packet.StopSound{SoundName: pk.SoundName, StopAll: pk.StopAll}}, true
	case *MobArmourEquipment:
		return []packet.Packet{&packet.MobArmourEquipment{EntityRuntimeID: pk.EntityRuntimeID, Helmet: pk.Helmet,
			Chestplate: pk.Chestplate, Leggings: pk.Leggings, Boots: pk.Boots}}, true
	case *InventoryContent:
		return []packet.Packet{&packet.InventoryContent{WindowID: pk.WindowID, Content: pk.Content}}, true
	case *InventorySlot:
		return []packet.Packet{&packet.InventorySlot{WindowID: pk.WindowID, Slot: pk.Slot, NewItem: pk.NewItem}}, true
	case *ChangeDimension, *Disconnect, *PlayerArmourDamage, *ResourcePacksInfo:
		return []packet.Packet{}, true
	}
	return nil, false
}
