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

package v712

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v729"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

func applyDeltas712(p packet.Pool) {
	delete(p, v786.IDCameraAimAssist)
	delete(p, v786.IDContainerRegistryCleanup)
	p[v786.IDUpdateAttributes] = func() packet.Packet { return &UpdateAttributes{} }
	p[v786.IDInventoryContent] = func() packet.Packet { return &InventoryContent{} }
	p[v786.IDInventorySlot] = func() packet.Packet { return &InventorySlot{} }
	p[v786.IDEmote] = func() packet.Packet { return &Emote{} }
	p[v786.IDTransfer] = func() packet.Packet { return &Transfer{} }
	p[v786.IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }
}

func NewClientPool() packet.Pool {
	p := v729.NewClientPool()
	applyDeltas712(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v729.NewServerPool()
	applyDeltas712(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v729.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v729.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.StartGame:
		out, _ := v729.FromLatestShared(proto, pk)
		sg := out[0].(*v766.StartGame)
		sg.BaseGameVersion, sg.GameVersion = "1.21.20", "1.21.20"
		items := itemdata.Items712()
		sg.Items = make([]v766.ItemEntry, len(items))
		for i, it := range items {
			sg.Items[i] = v766.ItemEntry{Name: it.Name, RuntimeID: it.RuntimeID, ComponentBased: it.ComponentBased}
		}
		return []packet.Packet{sg}, true
	case *packet.UpdateAttributes:
		return []packet.Packet{fromLatestUpdateAttributes(pk)}, true
	case *packet.InventoryContent:
		return []packet.Packet{fromLatestInventoryContent(pk)}, true
	case *packet.InventorySlot:
		return []packet.Packet{fromLatestInventorySlot(pk)}, true
	case *packet.Emote:
		return []packet.Packet{fromLatestEmote(pk)}, true
	case *packet.Transfer:
		return []packet.Packet{fromLatestTransfer(pk)}, true
	case *packet.ResourcePacksInfo:
		return []packet.Packet{fromLatestResourcePacksInfo(pk)}, true
	case *packet.CameraAimAssist, *packet.ContainerRegistryCleanup:
		return nil, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *UpdateAttributes:
		return []packet.Packet{toLatestUpdateAttributes(pk)}, true
	case *InventoryContent:
		return []packet.Packet{toLatestInventoryContent(pk)}, true
	case *InventorySlot:
		return []packet.Packet{toLatestInventorySlot(pk)}, true
	case *Emote:
		return []packet.Packet{toLatestEmote(pk)}, true
	case *Transfer:
		return []packet.Packet{v786.ToLatestTransfer(&v786.Transfer{Address: pk.Address, Port: pk.Port})}, true
	case *ResourcePacksInfo:
		return []packet.Packet{&packet.ResourcePacksInfo{}}, true
	}
	return nil, false
}
